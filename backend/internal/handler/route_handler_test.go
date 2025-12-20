package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/routing"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
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

	// set JWT key and generate a token for user 42
	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	// prepare user service with a mock user that has a cycling speed
	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp map[string]any
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	props, ok := resp["properties"].(map[string]any)
	require.True(t, ok)
	distVal, ok := props["distance_km"].(float64)
	require.True(t, ok)

	// expected distance should match routing.FindRouteFromCoordinates on the same coords
	expectedDist, _ := routing.FindRouteFromCoordinates(g, 0, 0, 1, 1, nil, nil)
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

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=10&endLongitude=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandleGetRoute_BadRequest(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(&models.Network{Nodes: map[int64]models.Node{}, Edges: map[int64][]models.Edge{}}, nil), userSvc)
	req, _ := http.NewRequest("GET", "/route?startLatitude=notanumber&startLongitude=0&endLatitude=0&endLongitude=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler := rh.HandleGetRoute()
	handler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestHandleGetRoute_GuestAccess - same as HappyPath but without authentication token
func TestHandleGetRoute_GuestAccess(t *testing.T) {
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

	utils.SetJWTKey([]byte("test-jwt-key"))

	// Mock user with cycling speed - but guest won't use it, will use default
	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 20}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/route?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	// NO Authorization header - guest access
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
	expectedDist, _ := routing.FindRouteFromCoordinates(g, 0, 0, 1, 1, nil, nil)
	require.InDelta(t, expectedDist, distVal, 1e-6)

	geom, ok := resp["geometry"].(map[string]any)
	require.True(t, ok)
	coords, ok := geom["coordinates"].([]any)
	require.True(t, ok)
	require.Len(t, coords, 3)
}

func TestHandleGetRoutes_HappyPath(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 1, Longitude: 0},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 4, Weight: 1.2}},
			2: {{To: 3, Weight: 1}},
			4: {{To: 3, Weight: 1.2}},
		},
	}

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1&k=2", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp map[string]any
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	require.Equal(t, "FeatureCollection", resp["type"])
	features, ok := resp["features"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, features)

	// Check first feature has correct structure
	firstFeature := features[0].(map[string]any)
	require.Equal(t, "Feature", firstFeature["type"])

	props, ok := firstFeature["properties"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, props, "route_index")
	require.Contains(t, props, "distance_km")
	require.Contains(t, props, "time_minutes")

	// Check route_index starts at 1
	routeIndex, ok := props["route_index"].(float64)
	require.True(t, ok)
	require.Equal(t, float64(1), routeIndex)

	geom, ok := firstFeature["geometry"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "LineString", geom["type"])
	coords, ok := geom["coordinates"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, coords)
}

func TestHandleGetRoutes_DefaultK(t *testing.T) {
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

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	// No k parameter - should default to 3
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]any
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	features, ok := resp["features"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, features)
}

func TestHandleGetRoutes_InvalidK(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(&models.Network{Nodes: map[int64]models.Node{}, Edges: map[int64][]models.Edge{}}, nil), userSvc)

	testCases := []struct {
		name string
		k    string
		code int
	}{
		{"negative k", "-1", http.StatusBadRequest},
		{"zero k", "0", http.StatusBadRequest},
		{"non-numeric k", "abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1&k="+tc.k, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()

			handler := rh.HandleGetRoutes()
			handler(rr, req)

			require.Equal(t, tc.code, rr.Code)
		})
	}
}

func TestHandleGetRoutes_Unreachable(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 10, Longitude: 10},
	}
	g := &models.Network{Nodes: nodes, Edges: map[int64][]models.Edge{1: {}, 2: {}}}

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=10&endLongitude=10&k=3", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandleGetRoutes_BadCoordinates(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(&models.Network{Nodes: map[int64]models.Node{}, Edges: map[int64][]models.Edge{}}, nil), userSvc)
	req, _ := http.NewRequest("GET", "/routes?startLatitude=notanumber&startLongitude=0&endLatitude=0&endLongitude=0&k=3", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleGetRoutes_GuestAccess(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 1, Longitude: 0},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 4, Weight: 1.2}},
			2: {{To: 3, Weight: 1}},
			4: {{To: 3, Weight: 1.2}},
		},
	}

	utils.SetJWTKey([]byte("test-jwt-key"))

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 20}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1&k=2", nil)
	// NO Authorization header - guest access
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp map[string]any
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	require.Equal(t, "FeatureCollection", resp["type"])
	features, ok := resp["features"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, features)

	// Verify all routes have time calculated with default speed
	for i, f := range features {
		feature := f.(map[string]any)
		props := feature["properties"].(map[string]any)
		distKm := props["distance_km"].(float64)
		timeMin := props["time_minutes"].(float64)

		// Time should be calculated using DefaultCyclingSpeed (15)
		expectedTime := distKm / float64(DefaultCyclingSpeed) * 60.0
		require.InDelta(t, expectedTime, timeMin, 1e-6, "route %d time mismatch", i+1)
	}
}

func TestHandleGetRoutes_MultipleRoutesIndexing(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 0, Longitude: 0},
		2: {ID: 2, Latitude: 0, Longitude: 1},
		3: {ID: 3, Latitude: 1, Longitude: 1},
		4: {ID: 4, Latitude: 1, Longitude: 0},
		5: {ID: 5, Latitude: 0.5, Longitude: 0.5},
	}
	g := &models.Network{
		Nodes: nodes,
		Edges: map[int64][]models.Edge{
			1: {{To: 2, Weight: 1}, {To: 4, Weight: 1.2}, {To: 5, Weight: 1.5}},
			2: {{To: 3, Weight: 1}},
			4: {{To: 3, Weight: 1.2}},
			5: {{To: 3, Weight: 1.3}},
		},
	}

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1&k=3", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]any
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	features, ok := resp["features"].([]any)
	require.True(t, ok)

	// Verify route indices are sequential starting from 1
	for i, f := range features {
		feature := f.(map[string]any)
		props := feature["properties"].(map[string]any)
		routeIndex := props["route_index"].(float64)
		require.Equal(t, float64(i+1), routeIndex, "route index should be %d", i+1)
	}
}

func TestHandleGetRoutes_FewerRoutesThanK(t *testing.T) {
	// Graph with only one possible route
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

	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 42, CyclingSpeed: 15}}
	userSvc := service.NewUserService(repo)

	rh := NewRouteHandler(service.NewRouteService(g, nil), userSvc)
	// Ask for 5 routes but only 1 exists
	req, _ := http.NewRequest("GET", "/routes?startLatitude=0&startLongitude=0&endLatitude=1&endLongitude=1&k=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := rh.HandleGetRoutes()
	handler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]any
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	features, ok := resp["features"].([]any)
	require.True(t, ok)
	// Should return only the routes that exist (less than k)
	require.GreaterOrEqual(t, len(features), 1)
	require.LessOrEqual(t, len(features), 5)
}
