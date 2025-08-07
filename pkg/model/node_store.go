package model

import (
	"database/sql"
	"errors"
	"log"
)

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

// GetNode retrieves a node from the database by its ID.
// Returns a pointer to the node if found, or an error if not found or query fails.
func (s *DBNodeStore) GetNode(id int64) (*DBNode, error) {
	n := &DBNode{}
	err := s.DB.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err == sql.ErrNoRows {
		return nil, errors.New("node not found")
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
