package model

type Network struct {
	Nodes map[int64]DBNode  // map node ids to nodes
	Edges map[int64][]DBWay // map node ids to lists of ways
}
