package model

type NodeService interface {
	Insert(n DBNode) error
	GetNode(id int64) (*DBNode, error)
	GetAllNodes() (map[int64]DBNode, error)
}
