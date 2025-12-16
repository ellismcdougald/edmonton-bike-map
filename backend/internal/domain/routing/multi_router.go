package routing

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// Route represents a path from start to end with its total distance
type Route struct {
	Distance float64 // total km distance of the route (estimated)
	Path     []int64 // sequence of node IDs
}

// FindMultipleRoutes finds up to k diverse routes between two nodes using a perturbation-based approach.
//
// Behavior:
//   - Runs Dijkstra repeatedly (up to k times), penalizing edges from previously found routes
//     to encourage different paths.
//   - The first route is always the shortest path; subsequent routes aim for diversity.
//   - Suppresses duplicate paths and returns routes ordered by distance/discovery.
//   - Returns an empty slice if k <= 0 and may stop early if no more paths exist.
func FindMultipleRoutes(network *models.Network, start, end int64, k int) []Route {
	if k <= 0 {
		return []Route{}
	}

	var routes []Route
	penaltyMultiplier := 1.5 // penalize previous route edges by 50%
	attempts := 0
	maxAttempts := k * 3 // try up to 3x requests to find k unique routes

	for len(routes) < k && attempts < maxAttempts {
		attempts++

		// Build edge weight overrides for previous routes
		edgeWeights := make(map[[2]int64]float64)
		for _, prevRoute := range routes {
			// Penalize each edge in the previous route
			for i := 0; i < len(prevRoute.Path)-1; i++ {
				fromNode := prevRoute.Path[i]
				toNode := prevRoute.Path[i+1]

				// Find original edge weight
				for _, edge := range network.Edges[fromNode] {
					if edge.To == toNode {
						// Apply penalty multiplier (increases with number of previous routes)
						multiplier := penaltyMultiplier * float64(len(routes)+1)
						edgeWeights[[2]int64{fromNode, toNode}] = edge.Weight * multiplier
						break
					}
				}
			}
		}

		// Find shortest path with penalized weights
		_, prev, found := dijkstra(network, start, end, edgeWeights)
		if !found {
			// No more paths available
			break
		}

		path := reconstructPath(prev, end)

		// Check if this path is a duplicate of any existing route
		isDuplicate := false
		for _, existing := range routes {
			if pathsEqual(path, existing.Path) {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			routeDistance := estimateDistance(path, network.Nodes)
			routes = append(routes, Route{Distance: routeDistance, Path: path})
		}
	}

	return routes
}

// FindMultipleRoutesFromCoordinates finds up to k diverse routes between the network nodes nearest the given start and end latitude/longitude coordinates.
// It maps each coordinate to the nearest node in the network and returns up to k distinct routes connecting those nodes.
func FindMultipleRoutesFromCoordinates(network *models.Network, startLatitude, startLongitude, endLatitude, endLongitude float64, k int) []Route {
	startNodeID := nearestNode(startLatitude, startLongitude, network.Nodes)
	endNodeID := nearestNode(endLatitude, endLongitude, network.Nodes)
	return FindMultipleRoutes(network, startNodeID, endNodeID, k)
}

// pathsEqual reports whether two paths contain the same sequence of node IDs.
// It returns true if both slices have the same length and identical elements in the same order.
func pathsEqual(path1, path2 []int64) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i := range path1 {
		if path1[i] != path2[i] {
			return false
		}
	}
	return true
}
