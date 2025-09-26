package model

import (
	"database/sql"
	"fmt"
	"log"
)

// TxNodeStore provides methods to interact with nodes using an existing transaction.
type TxNodeStore struct {
	Tx *sql.Tx
}

// Insert inserts a single DBNode within the transaction.
// Does nothing if the node ID already exists.
func (s *TxNodeStore) Insert(n DBNode) error {
	_, err := s.Tx.Exec(
		"INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		n.ID, n.Latitude, n.Longitude,
	)
	return err
}

// InsertBatch inserts multiple DBNodes within the transaction in a single multi-row statement.
// Does nothing for nodes that already exist.
func (s *TxNodeStore) InsertBatch(nodes []DBNode) error {
	if len(nodes) == 0 {
		return nil
	}

	query := "INSERT INTO nodes (id, latitude, longitude) VALUES "
	args := []interface{}{}
	for i, n := range nodes {
		if i > 0 {
			query += ", "
		}
		base := i*3 + 1
		query += fmt.Sprintf("($%d, $%d, $%d)", base, base+1, base+2)
		args = append(args, n.ID, n.Latitude, n.Longitude)
	}

	query += " ON CONFLICT (id) DO NOTHING"

	if _, err := s.Tx.Exec(query, args...); err != nil {
		return fmt.Errorf("insert batch nodes: %w", err)
	}

	return nil
}

// InsertBatchChunks inserts DBNodes in smaller batches to avoid exceeding parameter limits.
// Each batch is inserted using InsertBatch.
func (s *TxNodeStore) InsertBatchChunks(nodes []DBNode, batchSize int) error {
	if len(nodes) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 1000
	}

	for i := 0; i < len(nodes); i += batchSize {
		end := i + batchSize
		if end > len(nodes) {
			end = len(nodes)
		}
		batch := nodes[i:end]
		if err := s.InsertBatch(batch); err != nil {
			return fmt.Errorf("failed to insert batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

// GetNode retrieves a single node by ID using the transaction.
// Returns the DBNode or an error if not found.
func (s *TxNodeStore) GetNode(id int64) (*DBNode, error) {
	n := &DBNode{}
	err := s.Tx.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node %d not found", id)
		}
		return nil, err
	}
	return n, nil
}

// GetAllNodes retrieves all nodes within the transaction.
// Returns a map from node ID to DBNode.
func (s *TxNodeStore) GetAllNodes() (map[int64]DBNode, error) {
	rows, err := s.Tx.Query("SELECT id, latitude, longitude FROM nodes")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	nodes := make(map[int64]DBNode)
	for rows.Next() {
		var node DBNode
		if err := rows.Scan(&node.ID, &node.Latitude, &node.Longitude); err != nil {
			return nil, err
		}
		nodes[node.ID] = node
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
