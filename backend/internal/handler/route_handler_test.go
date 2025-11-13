package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/routing"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

func TestHandleGetRoute_HappyPath(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}},
			2: {{To: 3, Weight: 1}},
		},
	}

	rh := NewRouteHandler(service.NewRouteService(g))
	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp map[string]any
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	props, ok := resp["properties"].(map[string]any)
	require.True(t, ok)
	distVal, ok := props["distance_km"].(float64)
	require.True(t, ok)

	// expected distance should match routing.FindRouteFromCoordinates on the same coords
	expectedDist, _ := routing.FindRouteFromCoordinates(g, 0, 0, 1, 1)
	require.InDelta(t, expectedDist, distVal, 1e-6)

	geom, ok := resp["geometry"].(map[string]any)
	require.True(t, ok)
	coords, ok := geom["coordinates"].([]any)
	require.True(t, ok)
	require.Len(t, coords, 3)
}

func TestHandleGetRoute_Unreachable(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	g := &models.Network{Nodes: nodes, Edges: map[int64][]models.Edge{1: {}, 2: {}}}

	rh := NewRouteHandler(service.NewRouteService(g))
	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=10&endLongitude=10", nil)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandleGetRoute_BadRequest(t *testing.T) {
	rh := NewRouteHandler(service.NewRouteService(&models.Network{Nodes: map[int64]models.Node{}, Edges: map[int64][]models.Edge{}}))
	req, _ := http.NewRequest("GET", "/route?startLatitude=notanumber&startLongitude=0&endLatitude=0&endLongitude=0", nil)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
