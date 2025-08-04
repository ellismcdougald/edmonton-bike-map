package model

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/lib/pq"
)

// DBWay represents a way in the database, including its ID, tags as a map, and an ordered list of node IDs.
type DBWay struct {
	ID      int64
	Tags    map[string]string
	NodeIDs []int64
}

// Insert inserts the DBWay into the database, including its tags and the ordered node IDs.
// Uses a transaction to ensure atomicity.
// If the way ID already exists, the insert does nothing.
// Rolls back the transaction on any error.
func (w *DBWay) Insert(db *sql.DB) error {
	tagsJSON, err := json.Marshal(w.Tags)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
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

// GetAllWays retrieves all ways from the database along with their ordered node IDs.
// Returns a slice of DBWay structs or an error.
// Uses ARRAY_AGG to get node IDs in sequence order.
func GetAllWays(db *sql.DB) ([]DBWay, error) {
	query := `
		SELECT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags;
	`
	rows, err := db.Query(query)
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
