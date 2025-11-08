package models

// Edge represents a directed connection from one node to another with a weight
type Edge struct {
	To     int64   // ID of the destination node
	Weight float64 // computed weight (distance * multipliers)
}

// Network represents a graph made up of Nodes and Edges
type Network struct {
	Nodes map[int64]Node   // map node IDs to nodes
	Edges map[int64][]Edge // map node IDs to outgoing edges
}
