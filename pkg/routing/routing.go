package routing

import (
	"container/heap"
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

func FindRoute(network *model.Graph, start, end int64) (dist float64, path []int64) {
	distMap, prev, found := djikstra(network, start, end)
	if !found {
		return math.Inf(1), nil
	}
	path = reconstructPath(prev, end)
	dist = distMap[end]
	return dist, path
}

// dijkstra finds shortest path from start to goal and stops early when goal is reached.
// Returns:
// - dist: map[nodeID]distance from start to node
// - prev: map[nodeID]previous node in path for reconstruction
// - bool: true if goal reachable, false otherwise
func djikstra(g *model.Graph, start, goal int64) (dist map[int64]float64, prev map[int64]int64, found bool) {
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

// reconstructPath returns the shortest path from start to target using prev map.
func reconstructPath(prev map[int64]int64, target int64) []int64 {
	var path []int64
	for current, ok := target, true; ok; current, ok = prev[current] {
		path = append([]int64{current}, path...)
		if _, found := prev[current]; !found {
			break
		}
	}
	return path
}