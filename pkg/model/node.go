package model

import (
	"database/sql"
	"errors"
)

// DBNode represents a node with an ID and geographic coordinates stored in the database.
type DBNode struct {
	ID        int64
	Latitude  float64
	Longitude float64
}

// Insert inserts the DBNode into the database.
// If a node with the same ID already exists, the insert does nothing (ON CONFLICT DO NOTHING).
func (n *DBNode) Insert(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		n.ID, n.Latitude, n.Longitude,
	)
	return err
}

// GetNode retrieves a node from the database by its ID.
// Returns a pointer to the node if found, or an error if not found or query fails.
func GetNode(db *sql.DB, id int64) (*DBNode, error) {
	n := &DBNode{}
	err := db.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err == sql.ErrNoRows {
		return nil, errors.New("node not found")
	}
	return n, err
}

// GetAllNodes retrieves all nodes from the database.
// Returns a map from node ID to DBNode and any error encountered during the query.
func GetAllNodes(db *sql.DB) (map[int64]DBNode, error) {
	query := `
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
