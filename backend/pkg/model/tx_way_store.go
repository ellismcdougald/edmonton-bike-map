package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// TxWayStore provides methods to interact with ways and way_nodes using an existing transaction.
type TxWayStore struct {
	Tx *sql.Tx
}

// Insert inserts a single way and its nodes using the transaction.
// Does nothing if the way or its nodes already exist.
func (s *TxWayStore) Insert(w DBWay) error {
	tagsJSON, err := json.Marshal(w.Tags)
	if err != nil {
		return err
	}

	_, err = s.Tx.Exec(
		"INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		w.ID, tagsJSON,
	)
	if err != nil {
		return err
	}

	for seq, nodeID := range w.NodeIDs {
		_, err := s.Tx.Exec(
			"INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			w.ID, nodeID, seq,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertBatch inserts multiple ways and their nodes using the transaction.
// Uses a single multi-row INSERT for performance. Does nothing if the ways or nodes already exist.
func (s *TxWayStore) InsertBatch(ways []DBWay) error {
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

	if _, err := s.Tx.Exec(wayQuery, wayArgs...); err != nil {
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

	if _, err := s.Tx.Exec(nodeQuery, nodeArgs...); err != nil {
		return fmt.Errorf("insert batch way_nodes: %w", err)
	}

	return nil
}

// InsertBatchChunks splits a large list of ways into smaller chunks to avoid huge SQL statements.
// Each chunk is inserted using InsertBatch.
func (s *TxWayStore) InsertBatchChunks(ways []DBWay, batchSize int) error {
	if len(ways) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 500
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

// InsertBatchDynamic inserts ways in dynamically sized batches to avoid exceeding the maximum number of SQL parameters.
// Each batch is flushed once the parameter limit is reached.
func (s *TxWayStore) InsertBatchDynamic(ways []DBWay) error {
	const maxParams = 65000
	batch := []DBWay{}
	paramCount := 0

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := s.InsertBatch(batch); err != nil {
			return fmt.Errorf("insert dynamic batch: %w", err)
		}
		batch = nil
		paramCount = 0
		return nil
	}

	for _, w := range ways {
		// 2 params for ways + 3 params per node for way_nodes
		neededParams := 2 + 3*len(w.NodeIDs)
		if paramCount+neededParams > maxParams {
			if err := flush(); err != nil {
				return err
			}
		}
		batch = append(batch, w)
		paramCount += neededParams
	}

	return flush()
}
