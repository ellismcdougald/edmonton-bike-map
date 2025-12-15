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
	GetNearestWay(latitude, longitude float64) (*models.Way, error)
	GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error)
}

// ReviewRepository defines methods to interact with Review data.
type ReviewRepository interface {
	CreateReview(review *models.Review) error
	GetReviews(wayID int64) ([]models.Review, error)
	GetAllReviews() (map[int64][]models.Review, error)
	// InsertBatches inserts reviews in batches of the specified size.
	InsertBatches(reviews []models.Review, batchSize int) error
	// DeleteUserReviewForWay deletes the current user's review link for a specific way.
	// If the review is no longer linked to any ways after deletion, the review row may be removed.
	DeleteUserReviewForWay(userID int64, wayID int64) error
}

// UserRepository defines methods to interact with User data.
type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	Create(user *models.User) error
	UsernameExists(username string) (bool, error)
	UpdatePassword(userID int64, hashedPassword string) error
	UpdateCyclingSpeed(userID int64, speed int) error
}
