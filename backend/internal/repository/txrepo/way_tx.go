package txrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
	"github.com/lib/pq"
)

// TxWayRepository provides way operations within an existing transaction.
type TxWayRepository struct {
	Tx *sql.Tx
}

// NewTxWayRepository returns a WayRepository that operates using the provided transaction.
func NewTxWayRepository(tx *sql.Tx) repository.WayRepository {
	return &TxWayRepository{Tx: tx}
}

// Insert inserts a single way and its nodes using the transaction.
func (r *TxWayRepository) Insert(way models.Way) error {
	tagsJSON, err := json.Marshal(way.Tags)
	if err != nil {
		return err
	}

	_, err = r.Tx.Exec(
		"INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		way.ID, tagsJSON,
	)
	if err != nil {
		return err
	}

	for seq, nodeID := range way.NodeIDs {
		_, err := r.Tx.Exec(
			"INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			way.ID, nodeID, seq,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// insertBatch inserts multiple ways and associated way_nodes in a single multi-row statement.
func (r *TxWayRepository) insertBatch(ways []models.Way) error {
	if len(ways) == 0 {
		return nil
	}

	// --- Insert ways ---
	wayQuery := "INSERT INTO ways (id, tags) VALUES "
	wayArgs := []interface{}{}
	argIndex := 1
	for i, w := range ways {
		if i > 0 {
			wayQuery += ", "
		}
		tagsJSON, err := json.Marshal(w.Tags)
		if err != nil {
			return fmt.Errorf("marshal tags for way %d: %w", w.ID, err)
		}

		wayQuery += fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1)
		wayArgs = append(wayArgs, w.ID, tagsJSON)
		argIndex += 2
	}
	wayQuery += " ON CONFLICT (id) DO NOTHING"

	if _, err := r.Tx.Exec(wayQuery, wayArgs...); err != nil {
		return fmt.Errorf("insert batch ways: %w", err)
	}

	// --- Insert way_nodes ---
	nodeQuery := "INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES "
	nodeArgs := []interface{}{}
	argIndex = 1
	first := true
	for _, w := range ways {
		for seq, nodeID := range w.NodeIDs {
			if !first {
				nodeQuery += ", "
			}
			first = false
			nodeQuery += fmt.Sprintf("($%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2)
			nodeArgs = append(nodeArgs, w.ID, nodeID, seq)
			argIndex += 3
		}
	}
	nodeQuery += " ON CONFLICT DO NOTHING"

	if _, err := r.Tx.Exec(nodeQuery, nodeArgs...); err != nil {
		return fmt.Errorf("insert batch way_nodes: %w", err)
	}

	return nil
}

// InsertBatches inserts ways in batches of the specified size using the transaction.
func (r *TxWayRepository) InsertBatches(ways []models.Way, batchSize int) error {
	if batchSize <= 0 {
		return fmt.Errorf("invalid batch size: %d", batchSize)
	}

	for i := 0; i < len(ways); i += batchSize {
		end := i + batchSize
		if end > len(ways) {
			end = len(ways)
		}
		if err := r.insertBatch(ways[i:end]); err != nil {
			return fmt.Errorf("insert batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

// GetWay retrieves a way by its ID and returns the way with ordered node IDs.
func (r *TxWayRepository) GetWay(id int64) (*models.Way, error) {
	var way models.Way
	way.Tags = make(map[string]string)

	var tagsJson []byte

	err := r.Tx.QueryRow("SELECT id, tags FROM ways WHERE id = $1", id).
		Scan(&way.ID, &tagsJson)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrWayNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(tagsJson, &way.Tags); err != nil {
		return nil, err
	}

	rows, err := r.Tx.Query("SELECT node_id FROM way_nodes WHERE way_id = $1 ORDER BY sequence_id", id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var nodeID int64
		if err := rows.Scan(&nodeID); err != nil {
			return nil, err
		}
		way.NodeIDs = append(way.NodeIDs, nodeID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &way, nil
}

// GetAllWays retrieves all ways and their node IDs.
func (r *TxWayRepository) GetAllWays() ([]models.Way, error) {
	query := `
        SELECT
            w.id,
            w.tags,
            ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
        FROM ways w
        JOIN way_nodes wn ON wn.way_id = w.id
        GROUP BY w.id, w.tags;
    `
	rows, err := r.Tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	var ways []models.Way
	for rows.Next() {
		var way models.Way
		var tagsJson []byte

		err := rows.Scan(&way.ID, &tagsJson, pq.Array(&way.NodeIDs))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(tagsJson, &way.Tags)
		if err != nil {
			return nil, err
		}

		ways = append(ways, way)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ways, nil
}

// GetNearestWay retrieves the way whose geometry (constructed from ordered nodes) is closest to the given coordinates using the transaction.
// Uses PostGIS ST_Distance to find the closest point on the way line to the click point.
func (r *TxWayRepository) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	query := `
		WITH nearest_nodes AS (
			SELECT n.id
			FROM nodes n
			ORDER BY ST_SetSRID(ST_MakePoint(n.longitude, n.latitude), 4326)
					 <-> ST_SetSRID(ST_MakePoint($2, $1), 4326)
			LIMIT 400
		),
		candidate_ways AS (
			SELECT DISTINCT wn.way_id
			FROM way_nodes wn
			WHERE wn.node_id IN (SELECT id FROM nearest_nodes)
		),
		way_geometries AS (
			SELECT 
				w.id,
				w.tags,
				ST_MakeLine(
					ARRAY_AGG(
						ST_SetSRID(ST_MakePoint(n.longitude, n.latitude), 4326)
						ORDER BY wn.sequence_id
					)
				) AS geom
			FROM ways w
			JOIN way_nodes wn ON wn.way_id = w.id
			JOIN nodes n ON n.id = wn.node_id
			WHERE w.id IN (SELECT way_id FROM candidate_ways)
			GROUP BY w.id, w.tags
		),
		closest AS (
			SELECT 
				id,
				tags,
				ST_Distance(
					geom::geography,
					ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography
				)::double precision AS distance
			FROM way_geometries
			ORDER BY distance ASC
			LIMIT 1
		)
		SELECT 
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM closest c
		JOIN ways w ON w.id = c.id
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags
	`

	var way models.Way
	var tagsJson []byte

	err := r.Tx.QueryRow(query, latitude, longitude).Scan(&way.ID, &tagsJson, pq.Array(&way.NodeIDs))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrWayNotFound
		}
		return nil, err
	}

	way.Tags = make(map[string]string)
	err = json.Unmarshal(tagsJson, &way.Tags)
	if err != nil {
		return nil, err
	}

	return &way, nil
}

// GetWaysByNodeIDs retrieves all ways that contain any of the given node IDs using the transaction.
// Returns a slice of Ways that share at least one node with the provided node IDs.
func (r *TxWayRepository) GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error) {
	if len(nodeIDs) == 0 {
		return []models.Way{}, nil
	}

	query := `
		SELECT DISTINCT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		WHERE w.id IN (
			SELECT DISTINCT way_id 
			FROM way_nodes 
			WHERE node_id = ANY($1)
		)
		GROUP BY w.id, w.tags
	`

	rows, err := r.Tx.Query(query, pq.Array(nodeIDs))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var ways []models.Way
	for rows.Next() {
		var way models.Way
		var tagsJson []byte

		err := rows.Scan(&way.ID, &tagsJson, pq.Array(&way.NodeIDs))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(tagsJson, &way.Tags)
		if err != nil {
			return nil, err
		}

		ways = append(ways, way)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ways, nil
}
