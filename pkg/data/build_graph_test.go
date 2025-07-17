package data

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// BuildGraph:
func TestBuildGraphSimple(t *testing.T) {
	input := filepath.Join("testdata", "test_osm_data.json")

	network, err := BuildGraph(input)
	if err != nil {
		t.Fatalf("BuildGraph returned error: %v", err)
	}

	gotNumNodes := len(network.Nodes)
	wantNumNodes := 2
	if gotNumNodes != wantNumNodes {
		t.Errorf("Wrong node count: got %d, want %d", gotNumNodes, wantNumNodes)
	}

	gotNumEdges := len(network.Edges)
	wantNumEdges := 2
	if gotNumEdges != wantNumEdges {
		t.Errorf("Wrong edge count: got %d, want %d", gotNumEdges, wantNumEdges)
	}

	gotNodes := network.Nodes
	wantNodes := map[int64]model.Node{
		0: {Latitude: 53.2656877, Longitude: -113.6591141},
		1: {Latitude: 53.265261, Longitude: -113.6591135},
	}
	for id, wantNode := range wantNodes {
		gotNode, ok := gotNodes[id]
		if !ok {
			t.Errorf("Node with id %d not found in gotNodes", id)
		}
		if gotNode != wantNode {
			t.Errorf("Node mismatch for id %d: got %+v, want %+v", id, gotNode, wantNode)
		}
	}

	gotEdges := network.Edges
	wantEdges := map[int64][]model.Edge{
		0: {
			{To: 1, Weight: 47.44689197914839},
		},
		1: {
			{To: 0, Weight: 47.44689197914839},
		},
	}
	for id, wantNodeEdges := range wantEdges {
		gotNodeEdges, ok := gotEdges[id]
		if !ok {
			t.Errorf("Node with id %d not found in gotEdges", id)
		}

		if len(gotNodeEdges) != len(wantNodeEdges) {
			t.Errorf("Unequal edges for id %d: got %d, want %d", id, len(gotNodeEdges), len(wantNodeEdges))
			continue
		}
		for i, wantEdge := range wantNodeEdges {
			gotEdge := gotNodeEdges[i]
			if wantEdge != gotEdge {
				t.Errorf("Edge mismatch for id %d at index %d: got %+v, want %+v", id, i, gotEdge, wantEdge)
			}
		}
	}
}

// computeWayWeight:
func TestComputeWayWeightBikeLane(t *testing.T) {
	distance := 10.0
	tags := map[string]string{
		"bicycle":       "designated",
		"highway":       "cycleway",
		"lcn":           "yes",
		"motor_vehicle": "no",
	}

	wantWeight := distance * 0.5 * 1.75 * 0.9 * 0.9 * 0.7
	gotWeight := computeWayWeight(distance, tags)

	if math.Abs(gotWeight-wantWeight) > 1e-6 {
		t.Errorf("Wrong weight: got %f, wanted %f", gotWeight, wantWeight)
	}
}

func TestComputeWayWeightPoorBikeRoad(t *testing.T) {
	distance := 5.75
	tags := map[string]string{
		"highway": "trunk",
	}

	wantWeight := distance * math.Inf(1)
	gotWeight := computeWayWeight(distance, tags)

	if gotWeight != wantWeight {
		t.Errorf("Wrong weight: got %f, wanted %f", gotWeight, wantWeight)
	}
}
