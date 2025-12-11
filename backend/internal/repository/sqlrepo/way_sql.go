package sqlrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
	"github.com/lib/pq"
)

// SQLWayRepository provides methods to interact with the ways table in the database.
type SQLWayRepository struct {
	DB *sql.DB
}

// NewSQLWayRepository creates a new instance of SQLWayRepository.
func NewSQLWayRepository(db *sql.DB) repository.WayRepository {
	return &SQLWayRepository{DB: db}
}

// Insert adds a single way to the database.
// Uses a transaction to ensure atomicity.
// If the way ID already exists, the insert does nothing.
// Rolls back the transaction on any error.
func (r *SQLWayRepository) Insert(way models.Way) error {
	tagsJSON, err := json.Marshal(way.Tags)
	if err != nil {
		return err
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		way.ID, tagsJSON,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("Could not rollback transaction: %v", err)
		}
		return err
	}

	for seq, nodeID := range way.NodeIDs {
		_, err := tx.Exec(
			"INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			way.ID, nodeID, seq,
		)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("Could not rollback transaction: %v", err)
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// insertBatch adds multiple ways to the database in a single transaction.
// Inserts both the ways and their associated way_nodes.
// Uses ON CONFLICT DO NOTHING to skip existing entries.
// lower-case name: internal helper used by the exported InsertBatches method.
func (r *SQLWayRepository) insertBatch(ways []models.Way) error {
	if len(ways) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

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

	if _, err := tx.Exec(wayQuery, wayArgs...); err != nil {
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

	if _, err := tx.Exec(nodeQuery, nodeArgs...); err != nil {
		return fmt.Errorf("insert batch way_nodes: %w", err)
	}

	return tx.Commit()
}

// InsertBatches inserts ways in batches of the specified size.
// Uses the InsertBatch helper method.
func (r *SQLWayRepository) InsertBatches(ways []models.Way, batchSize int) error {
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

// GetWay retrieves a way by its ID.
// Also retrieves the associated node IDs in order.
func (r *SQLWayRepository) GetWay(id int64) (*models.Way, error) {
	var way models.Way
	way.Tags = make(map[string]string)
	var tagsJson []byte

	err := r.DB.QueryRow("SELECT id, tags FROM ways WHERE id = $1", id).
		Scan(&way.ID, &tagsJson)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("way with ID %d not found", id)
		}
		return nil, err
	}

	// Unmarshal tags JSON into the map
	if err := json.Unmarshal(tagsJson, &way.Tags); err != nil {
		return nil, err
	}

	rows, err := r.DB.Query("SELECT node_id FROM way_nodes WHERE way_id = $1 ORDER BY sequence_id", id)
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

// GetAllWays retrieves all ways from the database.
// Returns a slice of Ways and an error if the query fails.
func (r *SQLWayRepository) GetAllWays() ([]models.Way, error) {
	query := `
		SELECT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags;
	`
	rows, err := r.DB.Query(query)
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

// GetNearestWay retrieves the way whose nodes are closest to the given coordinates.
// Uses PostGIS ST_DWithin to find ways with nodes within a reasonable distance,
// then returns the way with the minimum distance to any of its nodes.
func (r *SQLWayRepository) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	query := `
		SELECT DISTINCT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		JOIN nodes n ON n.id = wn.node_id
		WHERE ST_DWithin(
			ST_SetSRID(ST_Point($1, $2), 4326)::geography,
			ST_SetSRID(ST_Point(n.longitude, n.latitude), 4326)::geography,
			5000  -- within 5 km
		)
		GROUP BY w.id, w.tags
		ORDER BY (
			SELECT MIN(ST_Distance(
				ST_SetSRID(ST_Point($1, $2), 4326)::geography,
				ST_SetSRID(ST_Point(n.longitude, n.latitude), 4326)::geography
			))
			FROM way_nodes wn2
			JOIN nodes n2 ON n2.id = wn2.node_id
			WHERE wn2.way_id = w.id
		) ASC
		LIMIT 1
	`

	var way models.Way
	var tagsJson []byte

	err := r.DB.QueryRow(query, longitude, latitude).Scan(&way.ID, &tagsJson, pq.Array(&way.NodeIDs))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no ways found near coordinates")
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
