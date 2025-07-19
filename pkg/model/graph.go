package model

// Node represents a location, consisting of a latitude and longitude.
type Node struct {
	Latitude, Longitude float64
}

// Edge represents a connection between nodes, with an associated weight.
type Edge struct {
	To     int64
	Weight float64
}

// Graph represents a graph structure consisting of nodes and edges between them.
type Graph struct {
	Nodes map[int64]Node   // map node ids to nodes
	Edges map[int64][]Edge // map node ids to lists of edges
}
