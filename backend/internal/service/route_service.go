package service

import (
	"fmt"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/routing"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

type RouteService struct {
	network *models.Network
}

// NewRouteService creates a new instance of RouteService.
func NewRouteService(network *models.Network) *RouteService {
	return &RouteService{
		network: network,
	}
}

// FindRoute returns the distance (km) and the list of nodes representing the route between two coordinates.
// If no route is found, distance will be -1 and nodes will be nil.
func (s *RouteService) FindRoute(startLatitude, startLongitude, endLatitude, endLongitude float64) (float64, []models.Node, error) {
	dist, ids := routing.FindRouteFromCoordinates(s.network, startLatitude, startLongitude, endLatitude, endLongitude)
	if dist < 0 || ids == nil {
		return dist, nil, nil
	}

	nodes := make([]models.Node, 0, len(ids))
	for _, id := range ids {
		n, ok := s.network.Nodes[id]
		if !ok {
			return 0, nil, fmt.Errorf("node id %d not found in network", id)
		}
		nodes = append(nodes, n)
	}

	return dist, nodes, nil
}

// FindMultipleRoutes returns the k shortest routes between two coordinates.
// Each route includes the distance (km) and the list of nodes representing the route.
// If no routes are found, an empty slice is returned.
func (s *RouteService) FindMultipleRoutes(startLatitude, startLongitude, endLatitude, endLongitude float64, k int) ([]RouteResult, error) {
	routes := routing.FindMultipleRoutesFromCoordinates(s.network, startLatitude, startLongitude, endLatitude, endLongitude, k)
	if len(routes) == 0 {
		return []RouteResult{}, nil
	}

	results := make([]RouteResult, 0, len(routes))
	for _, route := range routes {
		nodes := make([]models.Node, 0, len(route.Path))
		for _, id := range route.Path {
			n, ok := s.network.Nodes[id]
			if !ok {
				return nil, fmt.Errorf("node id %d not found in network", id)
			}
			nodes = append(nodes, n)
		}
		results = append(results, RouteResult{
			Distance: route.Distance,
			Nodes:    nodes,
		})
	}

	return results, nil
}

// RouteResult represents a single route with its distance and nodes
type RouteResult struct {
	Distance float64
	Nodes    []models.Node
}
