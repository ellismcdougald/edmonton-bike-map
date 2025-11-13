package service

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// NodeService provides operations related to nodes.
type NodeService struct {
	NodeRepository repository.NodeRepository
}

// NewNodeService creates a new instance of NodeService.
func NewNodeService(nodeRepo repository.NodeRepository) *NodeService {
	return &NodeService{
		NodeRepository: nodeRepo,
	}
}

// InsertNode inserts a new node using the NodeRepository.
func (s *NodeService) InsertNode(node models.Node) error {
	return s.NodeRepository.Insert(node)
}

// GetNodeByID retrieves a node by its ID using the NodeRepository.
func (s *NodeService) GetNodeByID(id int64) (*models.Node, error) {
	return s.NodeRepository.GetNode(id)
}

// GetAllNodes retrieves all nodes using the NodeRepository.
func (s *NodeService) GetAllNodes() (map[int64]models.Node, error) {
	return s.NodeRepository.GetAllNodes()
}
