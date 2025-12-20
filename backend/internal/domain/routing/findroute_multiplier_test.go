package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

// build a simple graph with two alternative paths between 1 and 4.
// Path A: 1->2->4 (wayID 10)
// Path B: 1->3->4 (wayID 20)
func makeTestNetwork() *models.Network {
	nodes := map[int64]models.Node{
		1: {ID: 1},
		2: {ID: 2},
		3: {ID: 3},
		4: {ID: 4},
	}

	edges := map[int64][]models.Edge{
		1: {
			{WayID: 10, To: 2, Weight: 1.0},
			{WayID: 20, To: 3, Weight: 1.25},
		},
		2: {
			{WayID: 10, To: 4, Weight: 1.0},
		},
		3: {
			{WayID: 20, To: 4, Weight: 1.25},
		},
	}

	return &models.Network{Nodes: nodes, Edges: edges}
}

func makeReviewsForRating(rating int, count int) []models.Review {
	out := make([]models.Review, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, models.Review{UserID: int64(i + 1), Rating: rating})
	}
	return out
}

func TestFindRoute_RespectsReviewMultipliers(t *testing.T) {
	g := makeTestNetwork()

	// Path A has many bad reviews -> multiplier > 1
	// Path B has many good reviews -> multiplier < 1
	reviews := map[int64][]models.Review{
		10: makeReviewsForRating(1, 10),
		20: makeReviewsForRating(10, 10),
	}
	mp := NewMultiplierProvider(reviews)

	_, path := findRoute(g, 1, 4, nil, mp)
	// expect path B chosen: 1 -> 3 -> 4
	require.Equal(t, []int64{1, 3, 4}, path)
}

func TestFindRoute_UserOverridePrefersTheirReview(t *testing.T) {
	g := makeTestNetwork()

	// Global reviews favour path B, but user has left a high rating on path A.
	reviews := map[int64][]models.Review{
		10: makeReviewsForRating(1, 10),
		20: makeReviewsForRating(10, 10),
	}
	// insert a user review (user id 999) with high rating for way 10
	userID := int64(999)
	userReview := models.Review{UserID: userID, Rating: 9}
	reviews[10] = append(reviews[10], userReview)

	mp := NewMultiplierProvider(reviews)

	_, path := findRoute(g, 1, 4, &userID, mp)
	// user-specific rating should favour path A: 1 -> 2 -> 4
	require.Equal(t, []int64{1, 2, 4}, path)
}
