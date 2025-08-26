package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lib/pq"
)

// DBWayStore provides methods to interact with the ways table in the database.
type DBWayStore struct {
	DB *sql.DB
}

// Insert inserts the DBWay into the database, including its tags and the ordered node IDs.
// Uses a transaction to ensure atomicity.
// If the way ID already exists, the insert does nothing.
// Rolls back the transaction on any error.
func (s *DBWayStore) Insert(w DBWay) error {
	tagsJSON, err := json.Marshal(w.Tags)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		w.ID, tagsJSON,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("Could not rollback transaction: %v", err)
		}
		return err
	}

	for seq, nodeID := range w.NodeIDs {
		_, err := tx.Exec(
			"INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			w.ID, nodeID, seq,
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

// InsertBatch inserts multiple DBWays in a single transaction.
// Inserts both the ways and their associated way_nodes.
// Uses ON CONFLICT DO NOTHING to skip existing entries.
func (s *DBWayStore) InsertBatch(ways []DBWay) error {
	if len(ways) == 0 {
		return nil
	}

	tx, err := s.DB.Begin()
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

// InsertBatchChunks inserts DBWays in smaller batches to avoid huge SQL statements.
// Each batch is inserted using the existing InsertBatch logic.
func (s *DBWayStore) InsertBatchChunks(ways []DBWay, batchSize int) error {
	if len(ways) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 500 // default batch size for ways
	}

	for i := 0; i < len(ways); i += batchSize {
		end := i + batchSize
		if end > len(ways) {
			end = len(ways)
		}
		batch := ways[i:end]
		if err := s.InsertBatch(batch); err != nil {
			return fmt.Errorf("failed to insert way batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

// GetAllWays retrieves all ways from the database along with their ordered node IDs.
// Returns a slice of DBWay structs or an error.
// Uses ARRAY_AGG to get node IDs in sequence order.
func (s *DBWayStore) GetAllWays() ([]DBWay, error) {
	query := `
		SELECT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var ways []DBWay
	for rows.Next() {
		var way DBWay
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
