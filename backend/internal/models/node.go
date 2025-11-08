package models

// Node represents a node with an ID and geographic coordinates stored in the database.
type Node struct {
	ID        int64
	Latitude  float64
	Longitude float64
}