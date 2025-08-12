package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
)

// MockNodeService implements model.NodeService
type MockNodeService struct {
	GetAllNodesFunc func() (map[int64]model.DBNode, error)
	GetNodeFunc     func(id int64) (*model.DBNode, error)
	InsertFunc      func(node model.DBNode) error
}

func (m *MockNodeService) GetAllNodes() (map[int64]model.DBNode, error) {
	if m.GetAllNodesFunc != nil {
		return m.GetAllNodesFunc()
	}
	return nil, errors.New("GetAllNodesFunc not implemented")
}

func (m *MockNodeService) GetNode(id int64) (*model.DBNode, error) {
	if m.GetNodeFunc != nil {
		return m.GetNodeFunc(id)
	}
	return nil, errors.New("GetNodeFunc not implemented")
}

func (m *MockNodeService) Insert(node model.DBNode) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(node)
	}
	return errors.New("InsertFunc not implemented")
}

// MockWayService implements model.WayService
type MockWayService struct {
	GetAllWaysFunc func() ([]model.DBWay, error)
	InsertFunc     func(way model.DBWay) error
}

func (m *MockWayService) GetAllWays() ([]model.DBWay, error) {
	if m.GetAllWaysFunc != nil {
		return m.GetAllWaysFunc()
	}
	return nil, errors.New("GetAllWaysFunc not implemented")
}

func (m *MockWayService) Insert(way model.DBWay) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(way)
	}
	return errors.New("InsertFunc not implemented")
}

// TestHandleAllWays tests the HandleAllWays handler for various scenarios including successful
// data retrieval, errors from NodeService and WayService, and missing nodes in ways.
func TestHandleAllWays(t *testing.T) {
	tests := []struct {
		name             string
		mockGetAllNodes  func() (map[int64]model.DBNode, error)
		mockGetAllWays   func() ([]model.DBWay, error)
		wantStatusCode   int
		wantFeatureCount int
		wantError        string // substring to look for in error response body
	}{
		{
			name: "Successful response",
			mockGetAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5, Longitude: -113.5},
					2: {ID: 2, Latitude: 53.6, Longitude: -113.6},
					3: {ID: 3, Latitude: 53.7, Longitude: -113.7},
				}, nil
			},
			mockGetAllWays: func() ([]model.DBWay, error) {
				return []model.DBWay{
					{ID: 10, NodeIDs: []int64{1, 2}, Tags: map[string]string{"name": "Way 1"}},
					{ID: 20, NodeIDs: []int64{2, 3}, Tags: map[string]string{"name": "Way 2"}},
				}, nil
			},
			wantStatusCode:   http.StatusOK,
			wantFeatureCount: 2,
		},
		{
			name: "NodeService GetAllNodes error",
			mockGetAllNodes: func() (map[int64]model.DBNode, error) {
				return nil, errors.New("db node error")
			},
			mockGetAllWays: func() ([]model.DBWay, error) {
				return nil, nil
			},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Could not get nodes from database",
		},
		{
			name: "WayService GetAllWays error",
			mockGetAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5, Longitude: -113.5},
				}, nil
			},
			mockGetAllWays: func() ([]model.DBWay, error) {
				return nil, errors.New("db way error")
			},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Could not get ways from database",
		},
		{
			name: "Missing node in way",
			mockGetAllNodes: func() (map[int64]model.DBNode, error) {
				return map[int64]model.DBNode{
					1: {ID: 1, Latitude: 53.5, Longitude: -113.5},
					// Node 2 missing
				}, nil
			},
			mockGetAllWays: func() ([]model.DBWay, error) {
				return []model.DBWay{
					{ID: 10, NodeIDs: []int64{1, 2}, Tags: map[string]string{"name": "Way 1"}},
				}, nil
			},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      "Node in way is missing from node data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNodeService := &MockNodeService{
				GetAllNodesFunc: tt.mockGetAllNodes,
			}
			mockWayService := &MockWayService{
				GetAllWaysFunc: tt.mockGetAllWays,
			}

			realHandlers := handlers.RealHandlers{
				NodeService: mockNodeService,
				WayService:  mockWayService,
			}

			handler := realHandlers.HandleAllWays()

			req := httptest.NewRequest(http.MethodGet, "/api/all-ways", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			resp := rec.Result()
			defer func() {
				if err := resp.Body.Close(); err != nil {
					// Log or handle error
					log.Printf("failed to close response body: %v", err)
				}
			}()

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("status code = %d; want %d", resp.StatusCode, tt.wantStatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if tt.wantStatusCode == http.StatusOK {
				var fc data.FeatureCollection
				if err := json.Unmarshal(body, &fc); err != nil {
					t.Fatalf("Failed to unmarshal response body as FeatureCollection: %v", err)
				}
				if len(fc.Features) != tt.wantFeatureCount {
					t.Errorf("feature count = %d; want %d", len(fc.Features), tt.wantFeatureCount)
				}
			} else {
				if !strings.Contains(string(body), tt.wantError) {
					t.Errorf("error response body = %q; want substring %q", string(body), tt.wantError)
				}
			}
		})
	}
}
