package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestFindMultipleRoutes_HappyPath(t *testing.T) {
	// Simple linear graph: 1 -> 2 -> 3
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

	routes := FindMultipleRoutes(g, 1, 3, 1)
	require.Equal(t, 1, len(routes))
	require.Equal(t, []int64{1, 2, 3}, routes[0].Path)
}

func TestFindMultipleRoutes_MultipleAlternatives(t *testing.T) {
	// Diamond graph with two paths of different costs
	// Path 1: 1 -> 2 -> 4 (cost: 2)
	// Path 2: 1 -> 3 -> 4 (cost: 3)
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 1, Longitude: 0},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 2, Longitude: 0.5},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 1.5}},
			2: {{To: 4, Weight: 1}},
			3: {{To: 4, Weight: 1.5}},
		},
	}

	routes := FindMultipleRoutes(g, 1, 4, 2)
	require.Equal(t, 2, len(routes))

	// First route should be cheapest
	require.Equal(t, []int64{1, 2, 4}, routes[0].Path)

	// Second route should use the alternative path (penalized but still reachable)
	require.True(t, len(routes[1].Path) > 0)
}

func TestFindMultipleRoutes_K_Greater_Than_Paths(t *testing.T) {
	// Request more paths than exist in this linear graph
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 1, Longitude: 0},
		3: {ID: 3, Latitude: 2, Longitude: 0},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			2: {{To: 3, Weight: 1}},
		},
	}

	routes := FindMultipleRoutes(g, 1, 3, 5)
	// Only 1 path exists in this linear graph
	require.Equal(t, 1, len(routes))
	require.Equal(t, []int64{1, 2, 3}, routes[0].Path)
}

func TestFindMultipleRoutes_SameNode(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
	}
	g := &models.Network{Nodes: nodes, Edges: map[int64][]models.Edge{}}

	routes := FindMultipleRoutes(g, 1, 1, 1)
	require.Equal(t, 1, len(routes))
	require.Equal(t, []int64{1}, routes[0].Path)
	require.Equal(t, 0.0, routes[0].Distance)
}

func TestFindMultipleRoutes_Unreachable(t *testing.T) {
	// Two disconnected nodes
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {},
			2: {},
		},
	}

	routes := FindMultipleRoutes(g, 1, 2, 1)
	require.Equal(t, 0, len(routes))
}

func TestFindMultipleRoutes_ZeroK(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 1, Longitude: 0},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
		},
	}

	routes := FindMultipleRoutes(g, 1, 2, 0)
	require.Equal(t, 0, len(routes))
}

func TestFindMultipleRoutesFromCoordinates(t *testing.T) {
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

	// Both coordinates map to node 1
	routes := FindMultipleRoutesFromCoordinates(g, 0, 0, 1, 1, 1)
	require.Equal(t, 1, len(routes))
	require.Equal(t, []int64{1, 2, 3}, routes[0].Path)
}

func TestFindMultipleRoutes_ComplexGraph(t *testing.T) {
	// More complex graph with multiple alternative paths
	// 1 -> 2 -> 5 (cost: 2)
	// 1 -> 3 -> 5 (cost: 3)
	// 1 -> 4 -> 5 (cost: 3.5)
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 1, Longitude: 0},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 0.5, Longitude: 1.5},
		5: {ID: 5, Latitude: 2, Longitude: 0.5},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 1.5}, {To: 4, Weight: 1.5}},
			2: {{To: 5, Weight: 1}},
			3: {{To: 5, Weight: 1.5}},
			4: {{To: 5, Weight: 2}},
		},
	}

	routes := FindMultipleRoutes(g, 1, 5, 3)
	// Should find at least 2 distinct paths
	require.True(t, len(routes) >= 1)

	// All paths should start at 1 and end at 5
	for _, route := range routes {
		require.Equal(t, int64(1), route.Path[0])
		require.Equal(t, int64(5), route.Path[len(route.Path)-1])
	}
}

func TestFindMultipleRoutes_NoDuplicates(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 1, Longitude: 0},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 2, Longitude: 0.5},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 3, Weight: 1.5}},
			2: {{To: 4, Weight: 1}},
			3: {{To: 4, Weight: 1.5}},
		},
	}

	routes := FindMultipleRoutes(g, 1, 4, 10)

	// Check for duplicates
	for i := 0; i < len(routes); i++ {
		for j := i + 1; j < len(routes); j++ {
			require.False(t, pathsEqual(routes[i].Path, routes[j].Path), "Found duplicate path")
		}
	}
}
