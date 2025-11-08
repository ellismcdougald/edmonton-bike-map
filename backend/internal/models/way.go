package models

// Way represents a path of nodes with an ID, tags as a map, and an ordered list of node IDs.
type Way struct {
	ID      int64
	Tags    map[string]string
	NodeIDs []int64
}
