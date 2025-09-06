package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// DBNodeStore provides methods to interact with the nodes table in the database.
type DBNodeStore struct {
	DB *sql.DB
}

// Insert inserts the DBNode into the database.
// If a node with the same ID already exists, the insert does nothing (ON CONFLICT DO NOTHING).
func (s *DBNodeStore) Insert(n DBNode) error {
	_, err := s.DB.Exec(
		"INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		n.ID, n.Latitude, n.Longitude,
	)
	return err
}

// InsertBatch inserts multiple DBNodes in a single query.
// If a node with the same ID already exists, it does nothing (ON CONFLICT DO NOTHING).
func (s *DBNodeStore) InsertBatch(nodes []DBNode) error {
	if len(nodes) == 0 {
		return nil
	}

	// Start a transaction for better performance
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	// Build the multi-row insert statement
	query := "INSERT INTO nodes (id, latitude, longitude) VALUES "
	args := []interface{}{}
	for i, n := range nodes {
		if i > 0 {
			query += ", "
		}
		// $1, $2, $3, ... for each node
		base := i*3 + 1
		query += fmt.Sprintf("($%d, $%d, $%d)", base, base+1, base+2)
		args = append(args, n.ID, n.Latitude, n.Longitude)
	}

	query += " ON CONFLICT (id) DO NOTHING"

	// Execute the batch insert
	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

// InsertBatchChunks inserts DBNodes in smaller batches to avoid huge SQL statements.
// Each batch is inserted using the existing InsertBatch logic.
func (s *DBNodeStore) InsertBatchChunks(nodes []DBNode, batchSize int) error {
	if len(nodes) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 1000 // default batch size
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

// GetNode retrieves a node from the database by its ID.
// Returns a pointer to the node if found, or an error if not found or query fails.
func (s *DBNodeStore) GetNode(id int64) (*DBNode, error) {
	n := &DBNode{}
	err := s.DB.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("node not found")
		}
		return nil, err
	}
	return n, err
}

// GetAllNodes retrieves all nodes from the database.
// Returns a map from node ID to DBNode and any error encountered during the query.
func (s *DBNodeStore) GetAllNodes() (map[int64]DBNode, error) {
	query := `
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;
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
