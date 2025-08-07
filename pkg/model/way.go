package model

// DBWay represents a way in the database, including its ID, tags as a map, and an ordered list of node IDs.
type DBWay struct {
	ID      int64
	Tags    map[string]string
	NodeIDs []int64
}
