package model

type Node struct {
	Latitude, Longitude float64
}

type Edge struct {
	To int64
	Weight float64
}

type Graph struct {
	Nodes map[int64]Node 			// map node ids to nodes
	Edges map[int64][]Edge		// map node ids to lists of edges
}