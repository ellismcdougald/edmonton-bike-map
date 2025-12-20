package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/routing"
	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

// Build a small network similar to findroute tests but with lat/lon for coordinates
func makeHandlerTestNetwork() *models.Network {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 0},
		4: {ID: 4, Latitude: 1, Longitude: 1},
	}

	// Two alternative paths from 1->4: via 2 (way 10) and via 3 (way 20)
	edges := map[int64][]models.Edge{
		1: {{WayID: 10, To: 2, Weight: 1.0}, {WayID: 20, To: 3, Weight: 1.5}},
		2: {{WayID: 10, To: 4, Weight: 1.0}},
		3: {{WayID: 20, To: 4, Weight: 1.0}},
	}

	return &models.Network{Nodes: nodes, Edges: edges}
}

func TestHandleGetRoute_UsesReviewMultipliers_Guest(t *testing.T) {
	g := makeHandlerTestNetwork()

	// reviews favour way 20 -> path 1->3->4 should be chosen
	reviews := map[int64][]models.Review{
		10: {{Rating: 1}, {Rating: 2}},
		20: {{Rating: 10}, {Rating: 9}},
	}
	mp := routing.NewMultiplierProvider(reviews)

	// provide a mock user service (not used for guest path selection)
	repo := &mockUserRepo{user: &models.User{ID: 1, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, mp), userSvc)

	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	geom := resp["geometry"].(map[string]any)
	coords := geom["coordinates"].([]any)

	// coordinates should reflect nodes 1 -> 3 -> 4
	// check length and final coordinate
	if len(coords) != 3 {
		t.Fatalf("expected 3 coords, got %d", len(coords))
	}
	last := coords[2].([]any)
	if last[0].(float64) != 1 || last[1].(float64) != 1 {
		t.Fatalf("unexpected final coord: %v", last)
	}
}

func TestHandleGetRoute_UsesReviewMultipliers_UserOverride(t *testing.T) {
	g := makeHandlerTestNetwork()

	reviews := map[int64][]models.Review{
		10: {{Rating: 1}, {Rating: 2}},
		20: {{Rating: 10}, {Rating: 9}},
	}
	// user 999 has left a high rating on way 10, so they should prefer path via 2
	userID := int64(999)
	reviews[10] = append(reviews[10], models.Review{UserID: userID, Rating: 9})

	mp := routing.NewMultiplierProvider(reviews)
	repo := &mockUserRepo{user: &models.User{ID: userID, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)
	rh := NewRouteHandler(service.NewRouteService(g, mp), userSvc)

	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	// inject user id into context (simulates optional auth middleware)
	req = req.WithContext(middleware.UserIDToContext(req.Context(), userID))
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	geom := resp["geometry"].(map[string]any)
	coords := geom["coordinates"].([]any)

	// coordinates should reflect nodes 1 -> 2 -> 4
	if len(coords) != 3 {
		t.Fatalf("expected 3 coords, got %d", len(coords))
	}
	second := coords[1].([]any)
	if second[0].(float64) != 1 || second[1].(float64) != 0 {
		t.Fatalf("unexpected second coord: %v", second)
	}
}
