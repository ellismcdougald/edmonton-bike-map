package routing

import (
	"container/heap"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// Route represents a path from start to end with its total distance
type Route struct {
	Distance float64 // total km distance of the route (estimated)
	Path     []int64 // sequence of node IDs
	cost     float64 // internal cost used during routing (sum of edge weights)
}

// FindMultipleRoutes returns the k shortest paths between two nodes using Yen's algorithm.
// It uses Dijkstra's algorithm as the underlying shortest path algorithm.
// Returns a slice of routes sorted by distance (shortest first).
// Distance is calculated as the km distance of the final route.
func FindMultipleRoutes(network *models.Network, start, end int64, k int) []Route {
	if k <= 0 {
		return []Route{}
	}

	var routes []Route
	var candidates *routeHeap

	// Find the shortest path
	distances, prev, found := dijkstra(network, start, end)
	if !found {
		return []Route{}
	}

	shortestPath := reconstructPath(prev, end)
	shortestCost := distances[end]
	routes = append(routes, Route{Distance: estimateDistance(shortestPath, network.Nodes), Path: shortestPath, cost: shortestCost})

	if k == 1 {
		return routes
	}

	// Initialize candidate paths heap
	candidates = &routeHeap{}
	heap.Init(candidates)

	// For each intermediate node in the shortest path
	for i := 0; i < len(shortestPath)-1; i++ {
		// Find spur path alternatives
		spurNode := shortestPath[i]
		rootPath := shortestPath[:i+1]

		// Remove edges used in previous shortest paths at this root
		removedEdges := removeEdgesFromNetwork(network, routes, rootPath)

		// Find shortest path from spur node to end, avoiding root path
		spurDistances, spurPrev, spurFound := dijkstra(network, spurNode, end)
		restoreEdgesToNetwork(network, removedEdges)

		if spurFound {
			// Reconstruct full path
			spurPath := reconstructPath(spurPrev, end)
			fullPath := append([]int64{}, rootPath...)
			fullPath = append(fullPath, spurPath[1:]...) // skip spur node (already in root)

			// Calculate cost as sum of root path cost + spur path cost
			rootCost := distances[spurNode]
			spurCost := spurDistances[end]
			fullCost := rootCost + spurCost

			// Add to candidates if not already in routes
			if !pathExists(routes, fullPath) {
				heap.Push(candidates, &candidateRoute{cost: fullCost, path: fullPath})
			}
		}
	}

	// Extract routes from candidates until we have k or candidates is empty
	for len(routes) < k && candidates.Len() > 0 {
		candidate := heap.Pop(candidates).(*candidateRoute)
		routes = append(routes, Route{Distance: estimateDistance(candidate.path, network.Nodes), Path: candidate.path, cost: candidate.cost})

		// For the newly added route, find more candidates
		for i := 0; i < len(candidate.path)-1; i++ {
			spurNode := candidate.path[i]
			rootPath := candidate.path[:i+1]

			removedEdges := removeEdgesFromNetwork(network, routes, rootPath)
			pathDistances, spurPrev, spurFound := dijkstra(network, spurNode, end)
			restoreEdgesToNetwork(network, removedEdges)

			if spurFound {
				spurPath := reconstructPath(spurPrev, end)
				fullPath := append([]int64{}, rootPath...)
				fullPath = append(fullPath, spurPath[1:]...)

				rootCost := distances[spurNode]
				spurCost := pathDistances[end]
				fullCost := rootCost + spurCost

				if !pathExists(routes, fullPath) && !candidateExists(candidates, fullPath) {
					heap.Push(candidates, &candidateRoute{cost: fullCost, path: fullPath})
				}
			}
		}
	}

	return routes
}

// FindMultipleRoutesFromCoordinates returns k shortest paths between two latitude-longitude pairs using Yen's algorithm.
func FindMultipleRoutesFromCoordinates(network *models.Network, startLatitude, startLongitude, endLatitude, endLongitude float64, k int) []Route {
	startNodeID := nearestNode(startLatitude, startLongitude, network.Nodes)
	endNodeID := nearestNode(endLatitude, endLongitude, network.Nodes)
	return FindMultipleRoutes(network, startNodeID, endNodeID, k)
}

// removeEdgesFromNetwork removes edges that are used in the root path of previous shortest paths.
// It returns the removed edges so they can be restored later.
func removeEdgesFromNetwork(network *models.Network, routes []Route, rootPath []int64) map[int64][]models.Edge {
	removed := make(map[int64][]models.Edge)

	for _, route := range routes {
		// Check if this route starts with the root path
		if len(route.Path) < len(rootPath) {
			continue
		}
		matches := true
		for i := 0; i < len(rootPath); i++ {
			if route.Path[i] != rootPath[i] {
				matches = false
				break
			}
		}

		if matches && len(route.Path) > len(rootPath) {
			// Remove the edge from the last node in rootPath to the next node
			fromNode := rootPath[len(rootPath)-1]
			toNode := route.Path[len(rootPath)]

			// Save removed edges
			if _, ok := removed[fromNode]; !ok {
				removed[fromNode] = network.Edges[fromNode]
				// Create a new slice without the edge to toNode
				newEdges := make([]models.Edge, 0)
				for _, edge := range network.Edges[fromNode] {
					if edge.To != toNode {
						newEdges = append(newEdges, edge)
					}
				}
				network.Edges[fromNode] = newEdges
			}
		}
	}

	return removed
}

// restoreEdgesToNetwork restores edges that were previously removed.
func restoreEdgesToNetwork(network *models.Network, removed map[int64][]models.Edge) {
	for nodeID, edges := range removed {
		network.Edges[nodeID] = edges
	}
}

// pathExists checks if a path already exists in the routes slice.
func pathExists(routes []Route, path []int64) bool {
	for _, route := range routes {
		if pathsEqual(route.Path, path) {
			return true
		}
	}
	return false
}

// candidateExists checks if a path already exists in the candidates heap.
func candidateExists(candidates *routeHeap, path []int64) bool {
	for _, candidate := range *candidates {
		if pathsEqual(candidate.path, path) {
			return true
		}
	}
	return false
}

// pathsEqual checks if two paths are identical.
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

// candidateRoute represents a candidate path in Yen's algorithm
type candidateRoute struct {
	cost float64 // cost of the path (sum of edge weights)
	path []int64 // sequence of node IDs
}

// routeHeap is a min-heap of candidate routes sorted by cost
type routeHeap []*candidateRoute

func (h routeHeap) Len() int           { return len(h) }
func (h routeHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h routeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *routeHeap) Push(x interface{}) {
	*h = append(*h, x.(*candidateRoute))
}

func (h *routeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
