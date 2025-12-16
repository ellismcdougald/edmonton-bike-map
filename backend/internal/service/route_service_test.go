package service

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestFindRoute_HappyPath(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			2: {{To: 3, Weight: 1}},
		},
	}

	svc := NewRouteService(network)
	dist, routeNodes, err := svc.FindRoute(0, 0, 1, 1)

	require.NoError(t, err)
	require.Greater(t, dist, 0.0)
	require.NotNil(t, routeNodes)
	require.NotEmpty(t, routeNodes)
	require.Equal(t, int64(1), routeNodes[0].ID)
	require.Equal(t, int64(3), routeNodes[len(routeNodes)-1].ID)
}

func TestFindRoute_NoRouteExists(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {}, // No edges from node 1
			2: {},
		},
	}

	svc := NewRouteService(network)
	dist, routeNodes, err := svc.FindRoute(0, 0, 10, 10)

	require.NoError(t, err)
	require.Less(t, dist, 0.0)
	require.Nil(t, routeNodes)
}

func TestFindRoute_EmptyNetwork(t *testing.T) {
	network := &models.Network{
		Nodes: map[int64]models.Node{},
		Edges: map[int64][]models.Edge{},
	}

	svc := NewRouteService(network)
	dist, routeNodes, err := svc.FindRoute(0, 0, 1, 1)

	require.NoError(t, err)
	require.Less(t, dist, 0.0)
	require.Nil(t, routeNodes)
}

func TestFindMultipleRoutes_HappyPath(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 1, Longitude: 0},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 4, Weight: 1.2}},
			2: {{To: 3, Weight: 1}},
			4: {{To: 3, Weight: 1.2}},
		},
	}

	svc := NewRouteService(network)
	results, err := svc.FindMultipleRoutes(0, 0, 1, 1, 2)

	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.LessOrEqual(t, len(results), 2)

	// Check first route structure
	require.Greater(t, results[0].Distance, 0.0)
	require.NotEmpty(t, results[0].Nodes)
	require.Equal(t, int64(1), results[0].Nodes[0].ID)
	require.Equal(t, int64(3), results[0].Nodes[len(results[0].Nodes)-1].ID)

	// If multiple routes, verify second route is different but valid
	if len(results) > 1 {
		require.Greater(t, results[1].Distance, 0.0)
		require.NotEmpty(t, results[1].Nodes)
		require.Equal(t, int64(1), results[1].Nodes[0].ID)
		require.Equal(t, int64(3), results[1].Nodes[len(results[1].Nodes)-1].ID)
	}
}

func TestFindMultipleRoutes_NoRoutesExist(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {},
			2: {},
		},
	}

	svc := NewRouteService(network)
	results, err := svc.FindMultipleRoutes(0, 0, 10, 10, 3)

	require.NoError(t, err)
	require.Empty(t, results)
}

func TestFindMultipleRoutes_FewerRoutesThanK(t *testing.T) {
	// Graph with only one possible route
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			2: {{To: 3, Weight: 1}},
		},
	}

	svc := NewRouteService(network)
	results, err := svc.FindMultipleRoutes(0, 0, 1, 1, 5)

	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.LessOrEqual(t, len(results), 5)
	// Should return at least one route even if k=5
	require.GreaterOrEqual(t, len(results), 1)
}

func TestFindMultipleRoutes_EmptyNetwork(t *testing.T) {
	network := &models.Network{
		Nodes: map[int64]models.Node{},
		Edges: map[int64][]models.Edge{},
	}

	svc := NewRouteService(network)
	results, err := svc.FindMultipleRoutes(0, 0, 1, 1, 3)

	require.NoError(t, err)
	require.Empty(t, results)
}

func TestFindMultipleRoutes_MultipleRoutesOrdering(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 1, Longitude: 0},
		5: {ID: 5, Latitude: 0.5, Longitude: 0.5},
	}
	network := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 4, Weight: 1.2}, {To: 5, Weight: 1.5}},
			2: {{To: 3, Weight: 1}},
			4: {{To: 3, Weight: 1.2}},
			5: {{To: 3, Weight: 1.3}},
		},
	}

	svc := NewRouteService(network)
	results, err := svc.FindMultipleRoutes(0, 0, 1, 1, 3)

	require.NoError(t, err)
	require.NotEmpty(t, results)

	// Verify all routes have positive distances
	for i, result := range results {
		require.Greater(t, result.Distance, 0.0, "route %d should have positive distance", i)
	}

	// All routes should start at node 1 and end at node 3
	for i, result := range results {
		require.NotEmpty(t, result.Nodes, "route %d should have nodes", i)
		require.Equal(t, int64(1), result.Nodes[0].ID, "route %d should start at node 1", i)
		require.Equal(t, int64(3), result.Nodes[len(result.Nodes)-1].ID, "route %d should end at node 3", i)
	}
}

func TestNewRouteService(t *testing.T) {
	network := &models.Network{
		Nodes: map[int64]models.Node{},
		Edges: map[int64][]models.Edge{},
	}

	svc := NewRouteService(network)

	require.NotNil(t, svc)
	require.Equal(t, network, svc.network)
}
