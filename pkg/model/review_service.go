package model

type ReviewService interface {
	CreateReview(review *Review) error
	GetReviews(wayID int64) ([]Review, error)
}
