package routing

import "github.com/ellismcdougald/edmonton-bike-map/internal/models"

// MultiplierProvider provides weight multipliers for ways based on reviews.
type MultiplierProvider struct {
	reviews map[int64][]models.Review // Map of way IDs to their reviews
}

// NewMultiplierProvider creates a new MultiplierProvider instance.
func NewMultiplierProvider(reviews map[int64][]models.Review) *MultiplierProvider {
	return &MultiplierProvider{reviews: reviews}
}

// MultiplierFor returns the weight multiplier for a given way and user.
func (mp *MultiplierProvider) MultiplierFor(wayID int64, userID *int64) float64 {
	reviewMultiplier := 1.0
	if reviews, exists := mp.reviews[wayID]; exists {
		reviewMultiplier = computeReviewMultiplier(reviews, userID)
	}

	return reviewMultiplier
}
