package routing

import (
	"math"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// pathPlanarDistance sums planarDistance over the path using the provided nodes.
func pathPlanarDistance(path []int64, nodes map[int64]model.Node) float64 {
	if len(path) < 2 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(path); i++ {
		n1 := nodes[path[i-1]]
		n2 := nodes[path[i]]
		total += planarDistance(n1.Latitude, n1.Longitude, n2.Latitude, n2.Longitude)
	}
	return total
}

func buildTestGraph() *model.Graph {
	nodes := map[int64]model.Node{
		1: {Latitude: 0, Longitude: 0},
		2: {Latitude: 0, Longitude: 1},
		3: {Latitude: 1, Longitude: 1},
		4: {Latitude: 1, Longitude: 0},
	}

	edges := map[int64][]model.Edge{
		1: {{To: 2, Weight: 1}, {To: 4, Weight: 4}},
		2: {{To: 3, Weight: 2}},
		3: {{To: 4, Weight: 1}},
		4: {},
	}

	return &model.Graph{Nodes: nodes, Edges: edges}
}

func TestSquaredEucDistance(t *testing.T) {
	tests := []struct {
		lat1, lon1, lat2, lon2 float64
		want                   float64
	}{
		{0, 0, 3, 4, 25},
		{1, 1, 1, 1, 0},
		{-1, -1, 2, 2, 18},
	}

	for _, tt := range tests {
		got := squaredEucDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
		if got != tt.want {
			t.Errorf("squaredEucDistance(%f,%f,%f,%f) = %f; want %f",
				tt.lat1, tt.lon1, tt.lat2, tt.lon2, got, tt.want)
		}

		// Test symmetry
		gotReverse := squaredEucDistance(tt.lat2, tt.lon2, tt.lat1, tt.lon1)
		if gotReverse != tt.want {
			t.Errorf("squaredEucDistance symmetry failed for inputs (%f,%f) and (%f,%f): got %f, want %f",
				tt.lat2, tt.lon2, tt.lat1, tt.lon1, gotReverse, tt.want)
		}
	}
}

func TestNearestNode(t *testing.T) {
	nodes := map[int64]model.Node{
		1: {Latitude: 0, Longitude: 0},
		2: {Latitude: 0, Longitude: 1},
		3: {Latitude: 1, Longitude: 1},
		4: {Latitude: 1, Longitude: 0},
	}

	tests := []struct {
		lat, lon float64
		wantID   int64
	}{
		{0, 0, 1},
		{0, 0.9, 2},
		{1, 1, 3},
		{0.75, 0.25, 4}, // uniquely closer to node 4 than others
		{0.25, 0.75, 2}, // uniquely closer to node 2 than others
	}

	for _, tt := range tests {
		got := nearestNode(tt.lat, tt.lon, nodes)
		if got != tt.wantID {
			t.Errorf("nearestNode(%v, %v) = %d; want %d", tt.lat, tt.lon, got, tt.wantID)
		}
	}
}

func TestFindRoute(t *testing.T) {
	graph := buildTestGraph()

	tests := []struct {
		start, end int64
		wantPaths  [][]int64
		wantFound  bool
	}{
		{1, 4, [][]int64{{1, 4}, {1, 2, 3, 4}}, true},
		{4, 1, nil, false},
		{2, 2, [][]int64{{2}}, true},
	}

	for _, tt := range tests {
		dist, path := findRoute(graph, tt.start, tt.end)
		found := path != nil

		if found != tt.wantFound {
			t.Errorf("findRoute(%d, %d): found = %v; want %v", tt.start, tt.end, found, tt.wantFound)
		}

		if !found {
			if dist != -1 {
				t.Errorf("findRoute(%d, %d): dist = %f; want -1", tt.start, tt.end, dist)
			}
			continue
		}

		expectedDist := pathPlanarDistance(path, graph.Nodes)
		if math.Abs(dist-expectedDist) > 1e-9 {
			t.Errorf("findRoute(%d, %d): dist = %f; want %f", tt.start, tt.end, dist, expectedDist)
		}

		matches := false
		for _, wp := range tt.wantPaths {
			if len(wp) == len(path) {
				match := true
				for i := range wp {
					if wp[i] != path[i] {
						match = false
						break
					}
				}
				if match {
					matches = true
					break
				}
			}
		}
		if !matches {
			t.Errorf("findRoute(%d, %d): path = %v; want one of %v", tt.start, tt.end, path, tt.wantPaths)
		}
	}
}

func TestDijkstra(t *testing.T) {
	graph := buildTestGraph()

	tests := []struct {
		name      string
		start     int64
		goal      int64
		wantFound bool
		wantDist  float64
		// list of valid shortest paths
		validPaths [][]int64
	}{
		{
			name:      "path exists 1->4",
			start:     1,
			goal:      4,
			wantFound: true,
			wantDist:  4.0,
			validPaths: [][]int64{
				{1, 4},
				{1, 2, 3, 4},
			},
		},
		{
			name:       "no path 4->1",
			start:      4,
			goal:       1,
			wantFound:  false,
			wantDist:   math.Inf(1),
			validPaths: nil,
		},
		{
			name:      "start equals goal",
			start:     2,
			goal:      2,
			wantFound: true,
			wantDist:  0,
			validPaths: [][]int64{
				{2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist, prev, found := dijkstra(graph, tt.start, tt.goal)

			if found != tt.wantFound {
				t.Errorf("found = %v; want %v", found, tt.wantFound)
			}

			if found {
				gotDist := dist[tt.goal]
				if math.Abs(gotDist-tt.wantDist) > 1e-9 {
					t.Errorf("distance to %d = %f; want %f", tt.goal, gotDist, tt.wantDist)
				}

				gotPath := reconstructPath(prev, tt.goal)
				if len(gotPath) == 0 {
					t.Fatalf("reconstructPath returned empty path")
				}

				// Check start and end node correctness quickly
				if gotPath[0] != tt.start || gotPath[len(gotPath)-1] != tt.goal {
					t.Errorf("path start/end mismatch: got %v; want start %d and end %d", gotPath, tt.start, tt.goal)
				}

				// Check if gotPath matches any valid path
				matches := false
				for _, valid := range tt.validPaths {
					if len(valid) == len(gotPath) {
						match := true
						for i := range valid {
							if valid[i] != gotPath[i] {
								match = false
								break
							}
						}
						if match {
							matches = true
							break
						}
					}
				}

				if !matches {
					t.Errorf("path %v does not match any valid shortest path %v", gotPath, tt.validPaths)
				}
			}
		})
	}
}

func TestFindRouteFromCoordinates(t *testing.T) {
	graph := buildTestGraph()

	tests := []struct {
		startLat, startLon float64
		endLat, endLon     float64
		wantPaths          [][]int64
		wantFound          bool
	}{
		{
			startLat: 0, startLon: 0,
			endLat: 1, endLon: 0,
			wantPaths: [][]int64{
				{1, 4},
				{1, 2, 3, 4},
			},
			wantFound: true,
		},
		{
			startLat: 1, startLon: 0,
			endLat: 0, endLon: 0,
			wantPaths: nil,
			wantFound: false,
		},
		{
			startLat: 0, startLon: 1,
			endLat: 0, endLon: 1,
			wantPaths: [][]int64{
				{2},
			},
			wantFound: true,
		},
	}

	for _, tt := range tests {
		dist, path := FindRouteFromCoordinates(graph, tt.startLat, tt.startLon, tt.endLat, tt.endLon)
		found := path != nil

		if found != tt.wantFound {
			t.Errorf("FindRouteFromCoordinates(%v,%v -> %v,%v): found = %v; want %v",
				tt.startLat, tt.startLon, tt.endLat, tt.endLon, found, tt.wantFound)
		}

		if !found {
			if dist != -1 {
				t.Errorf("FindRouteFromCoordinates(%v,%v -> %v,%v): dist = %f; want -1",
					tt.startLat, tt.startLon, tt.endLat, tt.endLon, dist)
			}
			continue
		}

		expectedDist := pathPlanarDistance(path, graph.Nodes)
		if math.Abs(dist-expectedDist) > 1e-9 {
			t.Errorf("FindRouteFromCoordinates(%v,%v -> %v,%v): dist = %f; want %f",
				tt.startLat, tt.startLon, tt.endLat, tt.endLon, dist, expectedDist)
		}

		matches := false
		for _, wp := range tt.wantPaths {
			if len(wp) == len(path) {
				match := true
				for i := range wp {
					if wp[i] != path[i] {
						match = false
						break
					}
				}
				if match {
					matches = true
					break
				}
			}
		}

		if !matches {
			t.Errorf("FindRouteFromCoordinates(%v,%v -> %v,%v): path = %v; want one of %v",
				tt.startLat, tt.startLon, tt.endLat, tt.endLon, path, tt.wantPaths)
		}
	}
}
