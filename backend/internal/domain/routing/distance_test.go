package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestPlanarDistance(t *testing.T) {
	// zero distance
	require.Equal(t, 0.0, planarDistance(0, 0, 0, 0))

	// one degree latitude difference should be ~111.32 km
	expected := 111.32
	result := planarDistance(0, 0, 1, 0)
	require.InDelta(t, expected, result, 0.01)

	// diagonal (1 deg lat, 1 deg lon at equator)
	// dx at equator: cos(0) * metersPerDeg * 1 => metersPerDeg
	// so distance sqrt( metersPerDeg^2 + metersPerDeg^2 ) / 1000 = sqrt(2)*111.32
	expectedDiag := 111.32 * 1.0 * 1.41421356237
	resultDiag := planarDistance(0, 0, 1, 1)
	require.InDelta(t, expectedDiag, resultDiag, 0.02)
}

func TestEstimateDistance_TwoNodes(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
		2: {ID: 2, Latitude: 53.55, Longitude: -113.49},
	}

	path := []int64{1, 2}
	want := planarDistance(nodes[1].Latitude, nodes[1].Longitude, nodes[2].Latitude, nodes[2].Longitude)
	got := estimateDistance(path, nodes)
	require.InDelta(t, want, got, 1e-6)
}

func TestEstimateDistance_MultiSegment(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}

	// distance should be planarDistance(1->2) + planarDistance(2->3)
	want := planarDistance(0, 0, 0, 1) + planarDistance(0, 1, 1, 1)
	got := estimateDistance([]int64{1, 2, 3}, nodes)
	require.InDelta(t, want, got, 1e-6)
}

func TestEstimateDistance_ShortPathAndMissingNode(t *testing.T) {
	// short path (<2) -> zero
	require.Equal(t, 0.0, estimateDistance([]int64{}, map[int64]models.Node{}))
	require.Equal(t, 0.0, estimateDistance([]int64{1}, map[int64]models.Node{1: {ID: 1}}))

	// missing node should panic
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
	}
	// path references node 2 which is missing
	require.Panics(t, func() { estimateDistance([]int64{1, 2}, nodes) })
}
