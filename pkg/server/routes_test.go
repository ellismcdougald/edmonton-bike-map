package server_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
)

type StubHandlers struct{}

func (h *StubHandlers) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HandleLogin"))
	}
}

func (h *StubHandlers) HandleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HandleSignup"))
	}
}

func (h *StubHandlers) HandleRouteByCoordinates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HandleRouteByCoordinates"))
	}
}

// TestRegisterRoutes checks that expected routes are registered and respond.
func TestRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	stubHandlers := &StubHandlers{}
	stubHandlers.HandleRouteByCoordinates()
	server.RegisterRoutes(mux, &model.Graph{}, &sql.DB{}, stubHandlers)

	handlers.HandleAllWays = func(w http.ResponseWriter, r *http.Request, _ *sql.DB) {
		w.Write([]byte("HandleAllWays"))
	}
	handlers.HandleGetReviews = func(w http.ResponseWriter, r *http.Request, _ *sql.DB) {
		w.Write([]byte("HandleGetReviews"))
	}
	handlers.HandlePostReview = func(w http.ResponseWriter, r *http.Request, _ *sql.DB) {
		w.Write([]byte("HandlePostReview"))
	}

	tests := []struct {
		method         string
		path           string
		wantStatus     int
		wantBody       string
		allowedOrigin  string
		allowedMethods string
	}{
		// GET requests
		{"GET", "/api/all-ways", http.StatusOK, "HandleAllWays", "http://localhost:5173", "GET, OPTIONS"},
		{"GET", "/api/route?startLatitude=53.5461&startLongitude=-113.4938&endLatitude=53.5444&endLongitude=-113.4909", http.StatusOK, "HandleRouteByCoordinates", "http://localhost:5173", "GET, OPTIONS"},
		{"GET", "/api/reviews?wayID=12345", http.StatusOK, "HandleGetReviews", "http://localhost:5173", "GET, POST, OPTIONS"},

		// POST requests
		{"POST", "/api/reviews", http.StatusOK, "HandlePostReview", "http://localhost:5173", "GET, POST, OPTIONS"},
		{"POST", "/api/signup", http.StatusOK, "HandleSignup", "http://localhost:5173", "POST, OPTIONS"},
		{"POST", "/api/login", http.StatusOK, "HandleLogin", "http://localhost:5173", "POST, OPTIONS"},

		// OPTIONS requests
		{"OPTIONS", "/api/route", http.StatusNoContent, "", "http://localhost:5173", "GET, OPTIONS"},
		{"OPTIONS", "/api/all-ways", http.StatusNoContent, "", "http://localhost:5173", "GET, OPTIONS"},
		{"OPTIONS", "/api/reviews", http.StatusNoContent, "", "http://localhost:5173", "GET, POST, OPTIONS"},
		{"OPTIONS", "/api/signup", http.StatusNoContent, "", "http://localhost:5173", "POST, OPTIONS"},
		{"OPTIONS", "/api/login", http.StatusNoContent, "", "http://localhost:5173", "POST, OPTIONS"},

		// Invalid route
		{"GET", "/api/nonexistent", http.StatusNotFound, "404 page not found\n", "", ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("%s %s: got status %d, want %d", tt.method, tt.path, rec.Code, tt.wantStatus)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.allowedOrigin {
			t.Errorf("%s %s: missing or wrong CORS header: expected %q, got %q", tt.method, tt.path, tt.allowedOrigin, got)
		}
		// Only check allowedMethods header if non-empty expected value (skip for 404)
		if tt.allowedMethods != "" {
			if got := rec.Header().Get("Access-Control-Allow-Methods"); got != tt.allowedMethods {
				t.Errorf("%s %s: wrong Access-Control-Allow-Methods header: expected %q, got %q", tt.method, tt.path, tt.allowedMethods, got)
			}
		}
		if gotBody := rec.Body.String(); gotBody != tt.wantBody {
			t.Errorf("Unexpected response body: wanted %q, got %q", tt.wantBody, gotBody)
		}
	}
}
