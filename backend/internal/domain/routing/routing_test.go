package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestFindRouteFromCoordinates_HappyPath(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			2: {{To: 3, Weight: 1}},
		},
	}

	d, path := FindRouteFromCoordinates(g, 0, 0, 1, 1, nil, nil)
	require.Equal(t, []int64{1, 2, 3}, path)
	want := estimateDistance(path, nodes)
	require.InDelta(t, want, d, 1e-9)
}

func TestFindRouteFromCoordinates_SameNode(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
	}
	g := &models.Network{Nodes: nodes, Edges: map[int64][]models.Edge{}}

	d, path := FindRouteFromCoordinates(g, 53.5461, -113.4938, 53.5461, -113.4938, nil, nil)
	require.Equal(t, []int64{1}, path)
	require.Equal(t, 0.0, d)
}

func TestFindRouteFromCoordinates_Unreachable(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	g := &models.Network{Nodes: nodes, Edges: map[int64][]models.Edge{
		1: {},
		2: {},
	}}

	d, path := FindRouteFromCoordinates(g, 0, 0, 10, 10, nil, nil)
	require.Equal(t, -1.0, d)
	require.Nil(t, path)
}
