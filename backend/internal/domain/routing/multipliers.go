package routing

import (
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

func clampRating(r int) int {
	if r < 1 {
		return 1
	}
	if r > 10 {
		return 10
	}
	return r
}

func computeReviewMultiplier(reviews []models.Review, userID *int64) float64 {
	if len(reviews) == 0 {
		return 1.0
	}

	total := 0
	count := 0
	foundUser := false

	for _, review := range reviews {
		if userID != nil && review.UserID == *userID {
			foundUser = true
			r := clampRating(review.Rating)
			total = r
			count = 1
			break
		}
		r := clampRating(review.Rating)
		total += r
		count++
	}

	average := float64(total) / float64(count)
	score := (average - 1.0) / 9.0
	multiplier := 1.3 - 0.75*score

	// Consider number of reviews in strength of multiplier
	var confidence float64
	if foundUser {
		confidence = 1.0
	} else {
		confidence = math.Min(1.0, float64(count)/10.0)
	}

	return 1.0 + (multiplier-1.0)*confidence
}
