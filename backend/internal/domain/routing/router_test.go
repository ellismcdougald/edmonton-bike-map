package routing

import (
	"math"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestReconstructPath_Basic(t *testing.T) {
	prev := map[int64]int64{
		3: 2,
		2: 1,
	}
	path := reconstructPath(prev, 3)
	require.Equal(t, []int64{1, 2, 3}, path)
}

func TestDijkstraSimple(t *testing.T) {
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1},
			2: {ID: 2},
			3: {ID: 3},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1.5}, {To: 3, Weight: 10}},
			2: {{To: 3, Weight: 2.5}},
		},
	}

	dist, prev, found := dijkstra(g, 1, 3)
	require.True(t, found)
	require.InDelta(t, 4.0, dist[3], 1e-9)
	require.Equal(t, int64(2), prev[3])
}

func TestDijkstraUnreachable(t *testing.T) {
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1},
			2: {ID: 2},
			3: {ID: 3},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			// node 3 has no incoming edges and is unreachable
		},
	}

	dist, _, found := dijkstra(g, 1, 3)
	require.False(t, found)
	require.True(t, math.IsInf(dist[3], 1))
}

func TestFindRoute_EndToEnd(t *testing.T) {
	// small 3-node network with coordinates so estimateDistance can compute
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 10}},
			2: {{To: 3, Weight: 1}},
		},
	}

	d, path := findRoute(g, 1, 3)
	require.Equal(t, []int64{1, 2, 3}, path)
	// distance returned by findRoute is estimateDistance(path, nodes)
	want := estimateDistance(path, nodes)
	require.InDelta(t, want, d, 1e-9)
}

func TestFindRoute_UnreachableAndSameNode(t *testing.T) {
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1, Latitude: 0, Longitude: 0},
			2: {ID: 2, Latitude: 0, Longitude: 1},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
		},
	}

	// unreachable
	d, path := findRoute(g, 2, 1)
	require.Equal(t, -1.0, d)
	require.Nil(t, path)

	// same node
	d2, path2 := findRoute(g, 1, 1)
	require.Equal(t, 0.0, d2)
	require.Equal(t, []int64{1}, path2)
}

func TestDijkstra_EqualWeightPaths(t *testing.T) {
	// Two equally-weighted shortest paths: 1-2-4 and 1-3-4 (both weight 2)
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1},
			2: {ID: 2},
			3: {ID: 3},
			4: {ID: 4},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 1}},
			2: {{To: 4, Weight: 1}},
			3: {{To: 4, Weight: 1}},
		},
	}

	_, prev, found := dijkstra(g, 1, 4)
	require.True(t, found)
	// reconstructPath should return a valid shortest path of length 3
	path := reconstructPath(prev, 4)
	require.Len(t, path, 3)
	require.Equal(t, int64(1), path[0])
	require.Equal(t, int64(4), path[len(path)-1])
}

func TestDijkstra_WithCycle(t *testing.T) {
	// Graph with a cycle; shortest path should avoid infinite loops
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1},
			2: {ID: 2},
			3: {ID: 3},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 5}},
			2: {{To: 3, Weight: 1}, {To: 1, Weight: 1}}, // cycle edge
			3: {},
		},
	}

	dist, _, found := dijkstra(g, 1, 3)
	require.True(t, found)
	require.InDelta(t, 2.0, dist[3], 1e-9)
}

func TestDijkstra_SelfLoop(t *testing.T) {
	// Self-loop on start node should not break algorithm
	g := &models.Network{
		Nodes: map[int64]models.Node{
			1: {ID: 1},
			2: {ID: 2},
		},
		Edges: map[int64][]models.Edge{
			1: {{To: 1, Weight: 0.5}, {To: 2, Weight: 1}},
			2: {},
		},
	}

	dist, _, found := dijkstra(g, 1, 2)
	require.True(t, found)
	require.InDelta(t, 1.0, dist[2], 1e-9)
}
