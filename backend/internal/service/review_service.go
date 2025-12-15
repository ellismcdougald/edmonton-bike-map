package service

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// ReviewService provides operations related to reviews.
type ReviewService struct {
	ReviewRepository repository.ReviewRepository
}

// NewReviewService creates a new instance of ReviewService.
func NewReviewService(reviewRepo repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		ReviewRepository: reviewRepo,
	}
}

// AddReview adds a new review using the ReviewRepository.
func (s *ReviewService) AddReview(review models.Review) error {
	return s.ReviewRepository.CreateReview(&review)
}

// GetReviewsForWay retrieves reviews for a specific way using the ReviewRepository.
func (s *ReviewService) GetReviewsForWay(wayID int64) ([]models.Review, error) {
	return s.ReviewRepository.GetReviews(wayID)
}

// GetAllReviews retrieves all reviews using the ReviewRepository.
func (s *ReviewService) GetAllReviews() (map[int64][]models.Review, error) {
	return s.ReviewRepository.GetAllReviews()
}

// DeleteUserReviewForWay removes the authenticated user's review association with a given way.
func (s *ReviewService) DeleteUserReviewForWay(userID int64, wayID int64) error {
	return s.ReviewRepository.DeleteUserReviewForWay(userID, wayID)
}
