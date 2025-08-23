package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	allowedMethods := []string{"GET", "POST", "PUT"}
	middleware := corsMiddleware(allowedMethods...)

	// A dummy handler to check if next is called
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the dummy handler with the middleware
	handler := middleware(nextHandler)

	// Test OPTIONS request — should NOT call next and respond 204 No Content
	reqOptions := httptest.NewRequest(http.MethodOptions, "/", nil)
	wOptions := httptest.NewRecorder()
	handler.ServeHTTP(wOptions, reqOptions)

	if wOptions.Result().StatusCode != http.StatusNoContent {
		t.Errorf("OPTIONS request status = %d; want %d", wOptions.Result().StatusCode, http.StatusNoContent)
	}
	if called {
		t.Errorf("next handler should NOT be called for OPTIONS request")
	}

	// Test GET request — should call next and set CORS headers
	called = false
	reqGet := httptest.NewRequest(http.MethodGet, "/", nil)
	reqGet.Header.Set("Origin", "http://localhost:5173")
	wGet := httptest.NewRecorder()
	handler.ServeHTTP(wGet, reqGet)

	resp := wGet.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET request status = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	if !called {
		t.Errorf("next handler should be called for GET request")
	}

	// Check headers
	wantOrigin := "http://localhost:5173"
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != wantOrigin {
		t.Errorf("Access-Control-Allow-Origin = %q; want %q", got, wantOrigin)
	}

	wantMethods := strings.Join(allowedMethods, ", ")
	if got := resp.Header.Get("Access-Control-Allow-Methods"); got != wantMethods {
		t.Errorf("Access-Control-Allow-Methods = %q; want %q", got, wantMethods)
	}

	wantHeaders := "Content-Type"
	if got := resp.Header.Get("Access-Control-Allow-Headers"); got != wantHeaders {
		t.Errorf("Access-Control-Allow-Headers = %q; want %q", got, wantHeaders)
	}
}
