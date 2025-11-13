package network

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestComputeHighwayMultiplier_KnownAndDefault(t *testing.T) {
	// known tag
	v := computeHighwayMultiplier("cycleway")
	require.InDelta(t, 0.9, v, 1e-9)

	// unknown tag -> default 1.5
	v2 := computeHighwayMultiplier("something-else")
	require.InDelta(t, 1.5, v2, 1e-9)
}

func TestComputeBikeFriendlyMultiplier_And_TagsMultiplier(t *testing.T) {
	// cycleway applies 0.9
	tags := map[string]string{tagCycleway: "yes"}
	m := computeBikeFriendlyMultiplier(tags)
	require.InDelta(t, 0.9, m, 1e-9)

	// bicycle yes applies 0.95
	tags2 := map[string]string{tagBicycle: "yes"}
	m2 := computeBikeFriendlyMultiplier(tags2)
	require.InDelta(t, 0.95, m2, 1e-9)

	// both combine multiplicatively
	tags3 := map[string]string{tagCycleway: "yes", tagBicycle: "yes"}
	m3 := computeBikeFriendlyMultiplier(tags3)
	require.InDelta(t, 0.9*0.95, m3, 1e-9)

	// tags multiplier: highway penalizes but bike-friendly should override highway>1
	tags4 := map[string]string{"highway": "secondary", tagCycleway: "yes"}
	mult := computeTagsMultiplier(tags4)
	// bikeFriendly 0.9, highway would be 1.5 but should be forced to 1 -> result 0.9
	require.InDelta(t, 0.9, mult, 1e-9)
}

func TestComputeReviewMultiplier_Behaviour(t *testing.T) {
	// no reviews -> 1
	r := computeReviewMultiplier([]models.Review{})
	require.InDelta(t, 1.0, r, 1e-9)

	// one review rating 5 -> average 5 -> multiplier = 1.2 - 0.1*5 = 0.7
	// confidence = 0.1 -> final = 1 + (0.7-1)*0.1 = 0.97
	revs := []models.Review{{Rating: 5}}
	r2 := computeReviewMultiplier(revs)
	require.InDelta(t, 0.97, r2, 1e-9)

	// ten reviews -> confidence 1 -> final equals multiplier for avg=5 (0.7)
	ten := make([]models.Review, 10)
	for i := 0; i < 10; i++ {
		ten[i] = models.Review{Rating: 5}
	}
	r3 := computeReviewMultiplier(ten)
	require.InDelta(t, 0.7, r3, 1e-9)
}
