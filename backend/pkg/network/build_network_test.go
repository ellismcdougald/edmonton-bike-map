package network

import (
	"math"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

type mockNodeService struct {
	getAllNodes func() (map[int64]model.DBNode, error)
}

func (m *mockNodeService) Insert(n model.DBNode) error             { return nil }
func (m *mockNodeService) GetNode(id int64) (*model.DBNode, error) { return nil, nil }
func (m *mockNodeService) GetAllNodes() (map[int64]model.DBNode, error) {
	return m.getAllNodes()
}

type mockWayService struct {
	getAllWays func() ([]model.DBWay, error)
}

type mockReviewService struct {
	getAllReviews func() (map[int][]model.Review, error)
}

func (m *mockReviewService) CreateReview(*model.Review) error               { return nil }
func (m *mockReviewService) GetReviews(wayID int64) ([]model.Review, error) { return nil, nil }
func (m *mockReviewService) GetAllReviews() (map[int][]model.Review, error) {
	return m.getAllReviews()
}

func (m *mockWayService) Insert(w model.DBWay) error { return nil }
func (m *mockWayService) GetAllWays() ([]model.DBWay, error) {
	return m.getAllWays()
}

func TestBuildNetwork_Normal(t *testing.T) {
	mockNodes := &mockNodeService{
		getAllNodes: func() (map[int64]model.DBNode, error) {
			return map[int64]model.DBNode{
				1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
				2: {ID: 2, Latitude: 53.5462, Longitude: -113.4939},
			}, nil
		},
	}

	mockWays := &mockWayService{
		getAllWays: func() ([]model.DBWay, error) {
			return []model.DBWay{
				{
					ID:      100,
					NodeIDs: []int64{1, 2},
					Tags:    map[string]string{"highway": "residential"},
				},
			}, nil
		},
	}

	mockReviews := &mockReviewService{
		getAllReviews: func() (map[int][]model.Review, error) {
			return map[int][]model.Review{}, nil // no reviews for these tests
		},
	}

	graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.Nodes))
	}

	edgesFrom1 := graph.Edges[1]
	if len(edgesFrom1) != 1 || edgesFrom1[0].To != 2 {
		t.Errorf("expected edge 1->2, got %+v", edgesFrom1)
	}

	edgesFrom2 := graph.Edges[2]
	if len(edgesFrom2) != 1 || edgesFrom2[0].To != 1 {
		t.Errorf("expected edge 2->1, got %+v", edgesFrom2)
	}

	if edgesFrom1[0].Weight <= 0 {
		t.Errorf("expected positive weight, got %f", edgesFrom1[0].Weight)
	}
}

func TestBuildNetwork_EdgeCases(t *testing.T) {
	t.Run("missing node", func(t *testing.T) {
		mockNodes := &mockNodeService{
			getAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
				}, nil
			},
		}

		mockWays := &mockWayService{
			getAllWays: func() ([]model.DBWay, error) {
				return []model.DBWay{
					{
						ID:      101,
						NodeIDs: []int64{1, 999}, // 999 does not exist
						Tags:    map[string]string{"highway": "residential"},
					},
				}, nil
			},
		}

		mockReviews := &mockReviewService{
			getAllReviews: func() (map[int][]model.Review, error) {
				return map[int][]model.Review{}, nil // no reviews for these tests
			},
		}

		graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(graph.Edges[1]) != 0 {
			t.Errorf("expected no edges from node 1 because the other node is missing, got %+v", graph.Edges[1])
		}
	})

	t.Run("one-way street", func(t *testing.T) {
		mockNodes := &mockNodeService{
			getAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
					2: {ID: 2, Latitude: 53.5462, Longitude: -113.4939},
				}, nil
			},
		}

		mockWays := &mockWayService{
			getAllWays: func() ([]model.DBWay, error) {
				return []model.DBWay{
					{
						ID:      102,
						NodeIDs: []int64{1, 2},
						Tags:    map[string]string{"highway": "residential", "one_way": "yes"},
					},
				}, nil
			},
		}

		mockReviews := &mockReviewService{
			getAllReviews: func() (map[int][]model.Review, error) {
				return map[int][]model.Review{}, nil // no reviews for these tests
			},
		}

		graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(graph.Edges[1]) != 1 || graph.Edges[1][0].To != 2 {
			t.Errorf("expected edge 1->2, got %+v", graph.Edges[1])
		}
		if len(graph.Edges[2]) != 0 {
			t.Errorf("expected no edge 2->1 for one-way, got %+v", graph.Edges[2])
		}
	})

	t.Run("infinite weight for motorway", func(t *testing.T) {
		mockNodes := &mockNodeService{
			getAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
					2: {ID: 2, Latitude: 53.5462, Longitude: -113.4939},
				}, nil
			},
		}

		mockWays := &mockWayService{
			getAllWays: func() ([]model.DBWay, error) {
				return []model.DBWay{
					{
						ID:      103,
						NodeIDs: []int64{1, 2},
						Tags:    map[string]string{"highway": "motorway"},
					},
				}, nil
			},
		}

		mockReviews := &mockReviewService{
			getAllReviews: func() (map[int][]model.Review, error) {
				return map[int][]model.Review{}, nil // no reviews for these tests
			},
		}

		graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !math.IsInf(graph.Edges[1][0].Weight, 1) {
			t.Errorf("expected infinite weight for motorway, got %f", graph.Edges[1][0].Weight)
		}
	})
}

func TestBuildNetwork_BikeFriendlyMultipliers(t *testing.T) {
	mockNodes := &mockNodeService{
		getAllNodes: func() (map[int64]model.DBNode, error) {
			return map[int64]model.DBNode{
				1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
				2: {ID: 2, Latitude: 53.5462, Longitude: -113.4939},
			}, nil
		},
	}

	mockReviews := &mockReviewService{
		getAllReviews: func() (map[int][]model.Review, error) {
			return map[int][]model.Review{}, nil // no reviews for these tests
		},
	}

	cases := []struct {
		name     string
		tags     map[string]string
		expected float64
	}{
		{"cycleway reduces weight", map[string]string{"highway": "residential", "cycleway": "lane"}, 0.9},
		{"bicycle designated reduces weight", map[string]string{"highway": "residential", "bicycle": "designated"}, 0.9},
		{"LCN reduces weight slightly", map[string]string{"highway": "residential", "lcn": "yes"}, 0.95},
		{"motor vehicle no reduces weight", map[string]string{"highway": "residential", "motor_vehicle": "no"}, 0.9},
		{"residential with no bike tags is 1x", map[string]string{"highway": "residential"}, 1.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockWays := &mockWayService{
				getAllWays: func() ([]model.DBWay, error) {
					return []model.DBWay{
						{ID: 200, NodeIDs: []int64{1, 2}, Tags: tc.tags},
					}, nil
				},
			}

			graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			edge := graph.Edges[1][0]

			nodes, err := mockNodes.GetAllNodes()
			if err != nil {
				t.Fatalf("unexpected error getting nodes: %v", err)
			}

			baseDistance := haversineDistance(
				nodes[1].Latitude,
				nodes[1].Longitude,
				nodes[2].Latitude,
				nodes[2].Longitude,
			)

			gotMultiplier := edge.Weight / baseDistance
			if math.Abs(gotMultiplier-tc.expected) > 0.01 {
				t.Errorf("expected multiplier ~%f, got %f", tc.expected, gotMultiplier)
			}
		})
	}
}

func TestBuildNetwork_ReviewMultipliers(t *testing.T) {
	mockNodes := &mockNodeService{
		getAllNodes: func() (map[int64]model.DBNode, error) {
			return map[int64]model.DBNode{
				1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938},
				2: {ID: 2, Latitude: 53.5462, Longitude: -113.4939},
			}, nil
		},
	}

	cases := []struct {
		name    string
		reviews []model.Review
	}{
		{"no reviews", []model.Review{}},
		{"high ratings", []model.Review{{Rating: 5}, {Rating: 5}}},
		{"mixed ratings", []model.Review{{Rating: 3}, {Rating: 4}}},
		{"low ratings", []model.Review{{Rating: 1}, {Rating: 2}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockWays := &mockWayService{
				getAllWays: func() ([]model.DBWay, error) {
					return []model.DBWay{
						{ID: 300, NodeIDs: []int64{1, 2}, Tags: map[string]string{"highway": "residential"}},
					}, nil
				},
			}

			mockReviews := &mockReviewService{
				getAllReviews: func() (map[int][]model.Review, error) {
					return map[int][]model.Review{
						300: tc.reviews,
					}, nil
				},
			}

			graph, err := BuildNetwork(mockNodes, mockWays, mockReviews)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			edge := graph.Edges[1][0]

			nodes, err := mockNodes.GetAllNodes()
			if err != nil {
				t.Fatalf("unexpected error getting nodes: %v", err)
			}

			baseDistance := haversineDistance(
				nodes[1].Latitude,
				nodes[1].Longitude,
				nodes[2].Latitude,
				nodes[2].Longitude,
			)

			gotMultiplier := edge.Weight / baseDistance

			// Compute expected multiplier according to current formula
			expected := 1.0
			if len(tc.reviews) > 0 {
				total := 0
				for _, r := range tc.reviews {
					total += r.Rating
				}
				average := float64(total) / float64(len(tc.reviews))
				multiplier := 1.2 - 0.1*average
				confidence := math.Min(1.0, float64(len(tc.reviews))/10.0)
				expected = 1.0 + (multiplier-1.0)*confidence
			}

			if math.Abs(gotMultiplier-expected) > 0.001 {
				t.Errorf("expected review multiplier ~%f, got %f", expected, gotMultiplier)
			}
		})
	}
}
