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