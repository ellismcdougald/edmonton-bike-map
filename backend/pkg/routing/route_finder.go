// Package routing provides functions to find optimal bike routes using a graph of nodes and edges representing the bike network.
// Pathfinding is via Dijkstra's algorithm. Start and end locations can be given as coordinates and the nearest node in the network will be used.
package routing

import (
	"container/heap"
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// squaredEucDistance computes the squared Euclidean distance between two latitude-longitude pairs.
func squaredEucDistance(lat1, lon1, lat2, lon2 float64) float64 {
	distLat := lat2 - lat1
	distLon := lon2 - lon1
	return distLat*distLat + distLon*distLon
}

// FindRouteFromCoordinates returns the lowest-cost route between two latitude-longitude pairs.
func FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (dist float64, path []int64) {
	startNodeID := nearestNode(startLatitude, startLongitude, network.Nodes)
	endNodeID := nearestNode(endLatitude, endLongitude, network.Nodes)
	return findRoute(network, startNodeID, endNodeID)
}

// nearestNode returns the ID of the closest node in the network to the given latitude-longitude pair.
func nearestNode(latitude float64, longitude float64, nodes map[int64]model.Node) (nodeID int64) {
	smallestDistance := math.Inf(1)
	var smallestNodeID int64 = -1
	for id, node := range nodes {
		distance := squaredEucDistance(latitude, longitude, node.Latitude, node.Longitude)
		if distance < smallestDistance {
			smallestDistance = distance
			smallestNodeID = id
		}
	}
	return smallestNodeID
}

// findRoute returns the lowest-cost route between two nodes using Dijkstra's algorithm.
func findRoute(network *model.Graph, start, end int64) (dist float64, path []int64) {
	_, prev, found := dijkstra(network, start, end)
	if !found {
		return -1, nil
	}
	path = reconstructPath(prev, end)
	dist = estimateDistance(path, network.Nodes)
	return dist, path
}

// estimateDistance returns the estimated distance of the given path in kilometres
func estimateDistance(path []int64, nodes map[int64]model.Node) float64 {
	if len(path) < 2 {
		return 0
	}

	dist := 0.0
	for i := 1; i < len(path); i++ {
		node1, ok1 := nodes[path[i-1]]
		node2, ok2 := nodes[path[i]]
		if !ok1 || !ok2 {
			panic("estimateDistance: node ID not found in graph")
		}
		segDist := planarDistance(node1.Latitude, node1.Longitude, node2.Latitude, node2.Longitude)
		dist += segDist
	}

	return dist
}

// dijkstra finds the shortest path from start to goal and stops early when goal is reached.
// Returns:
// - dist: map[nodeID]distance from start to each node
// - prev: map[nodeID]previous node in path for reconstruction
// - found: true if goal reachable, false otherwise
func dijkstra(g *model.Graph, start, goal int64) (dist map[int64]float64, prev map[int64]int64, found bool) {
	dist = make(map[int64]float64)
	prev = make(map[int64]int64)

	// Initialize all distances to infinity
	for node := range g.Nodes {
		dist[node] = math.Inf(1)
	}
	dist[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{NodeID: start, Distance: 0})

	visited := make(map[int64]bool)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)
		u := current.NodeID

		// Stop early if we reached the goal
		if u == goal {
			found = true
			break
		}

		if visited[u] {
			continue
		}
		visited[u] = true

		for _, edge := range g.Edges[u] {
			v := edge.To
			alt := dist[u] + edge.Weight
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				heap.Push(pq, &Item{NodeID: v, Distance: alt})
			}
		}
	}

	return dist, prev, found
}

// reconstructPath returns the shortest path from start to target using the prev map.
func reconstructPath(prev map[int64]int64, target int64) []int64 {
	var path []int64
	for current := target; ; {
		path = append(path, current)
		prevNode, ok := prev[current]
		if !ok {
			break
		}
		current = prevNode
	}
	// Reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// planarDistance returns distance in km
func planarDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const metersPerDeg = 111320.0 // approximate meters per degree
	latAvg := (lat1 + lat2) / 2 * math.Pi / 180

	dx := (lon2 - lon1) * math.Cos(latAvg) * metersPerDeg
	dy := (lat2 - lat1) * metersPerDeg

	return math.Sqrt(dx*dx+dy*dy) / 1000.0 // convert to km
}
