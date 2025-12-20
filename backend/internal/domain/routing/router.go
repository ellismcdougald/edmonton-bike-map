package routing

import (
	"container/heap"
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// findRoute finds the lowest-cost route between two nodes in the network.
// It returns the total estimated distance along the route and the sequence of node IDs from start to end.
// If no route exists between the nodes, it returns -1 and nil.
func findRoute(network *models.Network, start, end int64, userId *int64, mp *MultiplierProvider) (dist float64, path []int64) {
	_, prev, found := dijkstra(network, start, end, userId, nil, mp)
	if !found {
		return -1, nil
	}
	path = reconstructPath(prev, end)
	dist = estimateDistance(path, network.Nodes)
	return dist, path
}

// dijkstra computes shortest paths from `start` to `goal`, stopping early when `goal` is reached.
//
// Parameters:
//   - g: network graph containing nodes and adjacency lists
//   - start, goal: node IDs defining the source and target
//   - edgeWeights: optional overrides keyed by [2]int64{u, v} to replace the weight of edge u→v
//     (useful for penalizing edges from previously found routes)
//
// Returns:
// - dist: map of shortest distances from `start` to each node (math.Inf(1) for unreachable nodes)
// - prev: predecessor map for reconstructing a shortest path from `start`
// - found: true if `goal` was reached (finite distance), false otherwise
func dijkstra(g *models.Network, start, goal int64, userId *int64, edgeWeights map[[2]int64]float64, mp *MultiplierProvider) (dist map[int64]float64, prev map[int64]int64, found bool) {
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

			// Use overridden weight if provided, otherwise use edge weight
			weight := edge.Weight
			if edgeWeights != nil {
				if w, ok := edgeWeights[[2]int64{u, v}]; ok {
					weight = w
				}
			}
			if mp != nil {
				weight *= mp.MultiplierFor(edge.WayID, userId)
			}

			alt := dist[u] + weight
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
