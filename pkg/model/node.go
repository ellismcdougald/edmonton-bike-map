package model

import (
	"database/sql"
	"errors"
)

type DBNode struct {
	ID int64
	Latitude float64
	Longitude float64
}

func (n *DBNode) Insert(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING", n.ID, n.Latitude, n.Longitude)
	return err
}

func GetNode(db *sql.DB, id int64) (*DBNode, error) {
	n := &DBNode{}
	err := db.QueryRow("SELECT id, latitude, longitude FROM nodes where id = $1", id).Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err == sql.ErrNoRows {
		return nil, errors.New("node not found")
	}
	return n, err
}

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

		err := rows.Scan(&node.ID, &node.Latitude, &node.Longitude)
		if err != nil {
			return nil, err
		}

		nodes[node.ID] = node
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}