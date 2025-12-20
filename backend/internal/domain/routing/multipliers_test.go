package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestComputeReviewMultiplier_NoReviews(t *testing.T) {
	got := computeReviewMultiplier([]models.Review{}, nil)
	require.Equal(t, 1.0, got)
}

func TestComputeReviewMultiplier_SingleReview_NoUser(t *testing.T) {
	reviews := []models.Review{{UserID: 2, Rating: 10}}
	got := computeReviewMultiplier(reviews, nil)
	// expected: 0.955 (see formula in code)
	require.InDelta(t, 0.955, got, 1e-6)
}

func TestComputeReviewMultiplier_UserReviewOverrides(t *testing.T) {
	uid := int64(5)
	reviews := []models.Review{{UserID: 3, Rating: 2}, {UserID: 5, Rating: 8}, {UserID: 6, Rating: 9}}
	got := computeReviewMultiplier(reviews, &uid)
	// when user's review exists, multiplier uses that rating directly -> ~0.7166666667
	require.InDelta(t, 0.7166666666666667, got, 1e-9)
}

func TestComputeReviewMultiplier_MultipleReviews_ConfidenceScaling(t *testing.T) {
	reviews := []models.Review{{Rating: 5}, {Rating: 7}, {Rating: 9}}
	got := computeReviewMultiplier(reviews, nil)
	// average 7 -> multiplier 0.8, confidence 0.3 -> result 0.94
	require.InDelta(t, 0.94, got, 1e-9)
}

func TestComputeReviewMultiplier_ClampsRatings(t *testing.T) {
	// rating below 1 is clamped to 1; above 10 clamped to 10
	reviews := []models.Review{{Rating: 0}, {Rating: 11}}
	got := computeReviewMultiplier(reviews, nil)
	// after clamping both -> [1,10] average 5.5 -> score=(5.5-1)/9=0.5 -> multiplier=1.3-0.75*0.5=0.925
	// confidence = min(1,2/10)=0.2 -> result = 1 + (0.925-1)*0.2 = 1 - 0.015 = 0.985
	require.InDelta(t, 0.985, got, 1e-9)
}

func TestMultiplierProvider_MultiplierFor(t *testing.T) {
	reviewsMap := map[int64][]models.Review{
		42: {{Rating: 9}},
	}
	mp := NewMultiplierProvider(reviewsMap)

	// non-existent way -> 1.0
	require.Equal(t, 1.0, mp.MultiplierFor(1, nil))

	// existing way -> computed multiplier
	got := mp.MultiplierFor(42, nil)
	require.InDelta(t, 0.9633333333333334, got, 1e-9)
}
