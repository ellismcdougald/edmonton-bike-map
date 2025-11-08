package repository

import "github.com/ellismcdougald/edmonton-bike-map/internal/models"

// NodeRepository defines methods to interact with Node data.
type NodeRepository interface {
	Insert(node models.Node) error
	InsertBatches(nodes []models.Node, batchSize int) error
	GetNode(id int64) (*models.Node, error)
	GetAllNodes() (map[int64]models.Node, error)
}

// WayRepository defines methods to interact with Way data.
type WayRepository interface {
	Insert(way models.Way) error
	InsertBatches(ways []models.Way, batchSize int) error
	GetWay(id int64) (*models.Way, error)
	GetAllWays() ([]models.Way, error)
}

// ReviewRepository defines methods to interact with Review data.
type ReviewRepository interface {
	CreateReview(review *models.Review) error
	GetReviews(wayID int64) ([]models.Review, error)
	GetAllReviews() (map[int64][]models.Review, error)
}

