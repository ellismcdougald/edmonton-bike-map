package model

// ReviewService defines the interface for managing reviews in a data store.
// It supports creating a review and getting all reviews for a particular way.
type ReviewService interface {
	CreateReview(review *Review) error
	GetReviews(wayID int64) ([]Review, error)
}
