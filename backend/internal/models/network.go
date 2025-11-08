package models

// Network represents a graph made up of Nodes and Ways.
type Network struct {
	Nodes map[int64]Node   // map node ids to nodes
	Edges map[int64][]Way  // map node ids to lists of ways
}