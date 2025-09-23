// Package server serves the API endpoints for Edmonton Bike Map.
package server

import (
	"database/sql"
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
)

// RegisterRoutes registers HTTP handlers for all API endpoints on the given ServeMux.
// It applies CORS middleware with appropriate allowed methods to each route.
// Parameters:
//   - mux: the HTTP request multiplexer where routes are registered.
//   - network: the in-memory graph used by routing handlers.
//   - db: the database connection pool used by handlers requiring database access.
//
// Routes registered:
//   - GET /api/route       : compute bike routes between coordinates
//   - GET /api/all-ways    : fetch all ways from the database
//   - GET, POST /api/reviews : get or post reviews for ways
//   - POST /api/signup     : user signup
//   - POST /api/login      : user login
func RegisterRoutes(mux *http.ServeMux, network *model.Graph, db *sql.DB, handlerFuncs handlers.APIHandlers) {
	mux.Handle("/api/route", corsMiddleware("GET", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlerFuncs.HandleRouteByCoordinates()
		handler(w, r)
	})))

	mux.Handle("/api/routes", corsMiddleware("GET", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlerFuncs.HandleRoutesByCoordinates()
		handler(w, r)
	})))

	mux.Handle("/api/all-ways", corsMiddleware("GET", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlerFuncs.HandleAllWays()
		handler(w, r)
	})))

	mux.Handle("/api/reviews", corsMiddleware("GET", "POST", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			handler := handlerFuncs.HandleGetReviews()
			handler(w, r)
		case http.MethodPost:
			handler := handlerFuncs.HandlePostReview()
			handler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/signup", corsMiddleware("POST", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			handler := handlerFuncs.HandleSignup()
			handler(w, r)
		}
	})))

	mux.Handle("/api/login", corsMiddleware("POST", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			handler := handlerFuncs.HandleLogin()
			handler(w, r)
		}
	})))
}
