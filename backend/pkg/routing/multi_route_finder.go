package routing

import (
	"container/heap"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// CandidateItem stores a full path in the heap
type CandidateItem struct {
	Path     model.Path
	Distance float64
	index    int
}

// CandidatePQ is a min-heap of CandidateItems
type CandidatePQ []*CandidateItem

func (pq CandidatePQ) Len() int           { return len(pq) }
func (pq CandidatePQ) Less(i, j int) bool { return pq[i].Distance < pq[j].Distance }
func (pq CandidatePQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *CandidatePQ) Push(x any)       { item := x.(*CandidateItem); item.index = len(*pq); *pq = append(*pq, item) }
func (pq *CandidatePQ) Pop() any          { old := *pq; n := len(old); item := old[n-1]; *pq = old[0 : n-1]; return item }

// FindRoutesFromCoordinates returns up to k fast alternative paths between two coordinates.
func FindRoutesFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64, k int) []model.Path {
	startNodeID := nearestNode(startLatitude, startLongitude, network.Nodes)
	endNodeID := nearestNode(endLatitude, endLongitude, network.Nodes)
	if startNodeID == -1 || endNodeID == -1 {
		return nil
	}

	return kFastRoutes(network, startNodeID, endNodeID, k)
}

// kFastRoutes finds up to k alternative paths using greedy A* with edge penalties.
func kFastRoutes(graph *model.Graph, start, goal int64, k int) []model.Path {
	var shortestPaths []model.Path
	candidates := &CandidatePQ{}
	heap.Init(candidates)

	penalties := map[struct{ from, to int64 }]float64{}

	// First path: plain A*
	firstPath := aStar(graph, start, goal, penalties)
	if len(firstPath) == 0 {
		return nil
	}
	shortestPaths = append(shortestPaths, model.Path{Nodes: firstPath, Cost: pathCost(graph, firstPath, penalties)})

	for i := 1; i < k; i++ {
		for j := 0; j < len(shortestPaths[i-1].Nodes)-1; j++ {
			rootPath := shortestPaths[i-1].Nodes[:j+1]
			spurNode := shortestPaths[i-1].Nodes[j]

			// Add temporary penalties to previous paths sharing rootPath
			for _, p := range shortestPaths {
				if len(p.Nodes) > j && equalSlices(rootPath, p.Nodes[:j+1]) {
					from := p.Nodes[j]
					to := p.Nodes[j+1]
					penalties[struct{ from, to int64 }{from, to}] += 1e3
				}
			}

			// Compute spur path
			spurPath := aStar(graph, spurNode, goal, penalties)
			if len(spurPath) > 0 {
				totalNodes := append([]int64{}, rootPath[:len(rootPath)-1]...)
				totalNodes = append(totalNodes, spurPath...)
				totalCost := pathCost(graph, totalNodes, penalties)
				heap.Push(candidates, &CandidateItem{Path: model.Path{Nodes: totalNodes, Cost: totalCost}, Distance: totalCost})
			}

			// Remove temporary penalties
			for _, p := range shortestPaths {
				if len(p.Nodes) > j && equalSlices(rootPath, p.Nodes[:j+1]) {
					from := p.Nodes[j]
					to := p.Nodes[j+1]
					delete(penalties, struct{ from, to int64 }{from, to})
				}
			}
		}

		if candidates.Len() == 0 {
			break
		}
		next := heap.Pop(candidates).(*CandidateItem).Path
		shortestPaths = append(shortestPaths, next)
	}

	return shortestPaths
}

// aStar is a simple A* shortest-path finder using squaredEucDistance as heuristic and optional edge penalties.
func aStar(graph *model.Graph, start, goal int64, penalties map[struct{ from, to int64 }]float64) []int64 {
	dist := map[int64]float64{start: 0}
	prev := make(map[int64]int64)

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{NodeID: start, Distance: heuristic(graph.Nodes[start], graph.Nodes[goal])})

	visited := make(map[int64]bool)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)
		u := current.NodeID
		if visited[u] {
			continue
		}
		visited[u] = true

		if u == goal {
			break
		}

		for _, e := range graph.Edges[u] {
			v := e.To
			cost := e.Weight + penalties[struct{ from, to int64 }{u, v}]
			alt := dist[u] + cost
			if old, ok := dist[v]; !ok || alt < old {
				dist[v] = alt
				prev[v] = u
				heap.Push(pq, &Item{NodeID: v, Distance: alt + heuristic(graph.Nodes[v], graph.Nodes[goal])})
			}
		}
	}

	path := reconstructPath(prev, goal)
	return path
}

// heuristic returns squared Euclidean distance between two nodes.
func heuristic(a, b model.Node) float64 {
	return squaredEucDistance(a.Latitude, a.Longitude, b.Latitude, b.Longitude)
}

// pathCost computes total cost of a path including penalties
func pathCost(graph *model.Graph, nodes []int64, penalties map[struct{ from, to int64 }]float64) float64 {
	cost := 0.0
	for i := 0; i < len(nodes)-1; i++ {
		from := nodes[i]
		to := nodes[i+1]
		for _, e := range graph.Edges[from] {
			if e.To == to {
				cost += e.Weight
				break
			}
		}
		cost += penalties[struct{ from, to int64 }{from, to}]
	}
	return cost
}

// equalSlices compares two int64 slices
func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
