package model

import (
	"database/sql"
	"encoding/json"
)

type DBWay struct {
	ID	int64
	Tags map[string]string
	NodeIDs []int64
}

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
		tx.Rollback()
		return err
	}

	for seq, nodeID := range w.NodeIDs {
		_, err := tx.Exec(
			"INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			w.ID, nodeID, seq,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}