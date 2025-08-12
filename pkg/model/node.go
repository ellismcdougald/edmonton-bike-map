package model

// DBNode represents a node with an ID and geographic coordinates stored in the database.
type DBNode struct {
	ID        int64
	Latitude  float64
	Longitude float64
}
