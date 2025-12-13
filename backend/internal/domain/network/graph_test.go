package network

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestBuildGraph_BasicEdgesAndOnewayAndMissingNodes(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
		2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
	}

	// basic way (no tags)
	ways := []models.Way{{ID: 11, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}}
	reviews := map[int64][]models.Review{}

	net, err := buildGraph(nodes, ways, reviews)
	require.NoError(t, err)

	// edges from 1 -> 2 and 2 -> 1 (not oneway)
	require.Len(t, net.Edges[1], 1)
	require.Len(t, net.Edges[2], 1)

	// confirm weight equals haversine*tagsMultiplier (highway empty -> default 1.5)
	w := net.Edges[1][0].Weight
	expectedDist := haversineDistance(nodes[1].Latitude, nodes[1].Longitude, nodes[2].Latitude, nodes[2].Longitude)
	require.InDelta(t, expectedDist*computeTagsMultiplier(map[string]string{}), w, 1e-6)

	// oneway should prevent reverse edge
	ways2 := []models.Way{{ID: 12, Tags: map[string]string{"one_way": "yes"}, NodeIDs: []int64{1, 2}}}
	net2, err := buildGraph(nodes, ways2, reviews)
	require.NoError(t, err)
	require.Len(t, net2.Edges[1], 1)
	require.Len(t, net2.Edges[2], 0)

	// missing nodes are skipped
	ways3 := []models.Way{{ID: 13, Tags: map[string]string{}, NodeIDs: []int64{1, 3}}}
	net3, err := buildGraph(nodes, ways3, reviews)
	require.NoError(t, err)
	// no edges should be added because 3 doesn't exist
	require.Empty(t, net3.Edges[1])
}

func TestBuildGraph_RespectsReviewMultiplier(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
		2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
	}

	ways := []models.Way{{ID: 21, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}}
	// one review with rating 5 -> review multiplier ~0.97
	reviews := map[int64][]models.Review{
		21: {{WayIDs: []int64{21}, Rating: 5}},
	}

	net, err := buildGraph(nodes, ways, reviews)
	require.NoError(t, err)

	require.Len(t, net.Edges[1], 1)
	w := net.Edges[1][0].Weight
	expectedDist := haversineDistance(nodes[1].Latitude, nodes[1].Longitude, nodes[2].Latitude, nodes[2].Longitude)
	expected := expectedDist * computeTagsMultiplier(map[string]string{}) * computeReviewMultiplier(reviews[21])
	require.InDelta(t, expected, w, 1e-6)
}
