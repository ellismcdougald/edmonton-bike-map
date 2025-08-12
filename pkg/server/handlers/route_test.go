package handlers_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
)

type MockRouter struct{}

func (m *MockRouter) FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (float64, []int64) {
	return 10.5, []int64{1, 2, 3}
}

func TestHandleRouteByCoordinates(t *testing.T) {
	mockRouter := &MockRouter{}

	graph := &model.Graph{
		Nodes: map[int64]model.Node{
			1: {Latitude: 53.5, Longitude: -113.5},
			2: {Latitude: 53.6, Longitude: -113.6},
			3: {Latitude: 53.7, Longitude: -113.7},
		},
	}

	realHandlers := handlers.RealHandlers{
		Router:  mockRouter,
		Network: graph,
	}

	handler := realHandlers.HandleRouteByCoordinates()

	tests := []struct {
		name            string
		url             string
		wantStatus      int
		wantCoordinates [][2]float64
	}{
		{
			name:       "Valid request 1",
			url:        "/route?startLatitude=53.5&startLongitude=-113.5&endLatitude=53.6&endLongitude=-113.6",
			wantStatus: http.StatusOK,
			wantCoordinates: [][2]float64{
				{-113.5, 53.5},
				{-113.6, 53.6},
				{-113.7, 53.7},
			},
		},
		{
			name:       "Valid request 2",
			url:        "/route?startLatitude=53.6&startLongitude=-113.6&endLatitude=53.7&endLongitude=-113.7",
			wantStatus: http.StatusOK,
			wantCoordinates: [][2]float64{
				{-113.5, 53.5},
				{-113.6, 53.6},
				{-113.7, 53.7},
			},
		},
		{
			name:       "Missing parameter",
			url:        "/route?startLatitude=53.5&startLongitude=-113.5&endLatitude=53.6",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid parameter",
			url:        "/route?startLatitude=53.5&startLongitude=bad&endLatitude=53.6&endLongitude=-113.6",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			handler(w, req)

			resp := w.Result()
			defer func() {
				if err := resp.Body.Close(); err != nil {
					// Log or handle error
					log.Printf("failed to close response body: %v", err)
				}
			}()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read body: %v", err)
				}

				var geojson struct {
					Type     string `json:"type"`
					Geometry struct {
						Type        string       `json:"type"`
						Coordinates [][2]float64 `json:"coordinates"`
					} `json:"geometry"`
					Properties map[string]interface{} `json:"properties"`
				}

				if err := json.Unmarshal(body, &geojson); err != nil {
					t.Fatalf("failed to unmarshal json: %v", err)
				}

				if geojson.Type != "Feature" {
					t.Errorf("unexpected type: got %q, want %q", geojson.Type, "Feature")
				}
				if geojson.Geometry.Type != "LineString" {
					t.Errorf("unexpected geometry type: got %q, want %q", geojson.Geometry.Type, "LineString")
				}

				if len(geojson.Geometry.Coordinates) != len(tt.wantCoordinates) {
					t.Fatalf("unexpected number of coordinates: got %d, want %d", len(geojson.Geometry.Coordinates), len(tt.wantCoordinates))
				}

				for i, coord := range geojson.Geometry.Coordinates {
					wantCoord := tt.wantCoordinates[i]
					if coord[0] != wantCoord[0] || coord[1] != wantCoord[1] {
						t.Errorf("coordinate %d = %v; want %v", i, coord, wantCoord)
					}
				}
			}
		})
	}
}
