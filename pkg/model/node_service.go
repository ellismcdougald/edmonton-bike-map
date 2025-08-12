package model

// NodeService defines the interface for managing nodes in a data store.
// It supports inserting individual nodes, getting individual nodes, and getting all nodes.
type NodeService interface {
	Insert(n DBNode) error
	GetNode(id int64) (*DBNode, error)
	GetAllNodes() (map[int64]DBNode, error)
}

// Note: NodeService should ideally work with domain-level Node models rather than DBNode directly