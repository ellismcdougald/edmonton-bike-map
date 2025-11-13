package server

import (
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/internal/handler"
)

// RegisterRoutes registers HTTP handlers for all API endpoints on the given ServeMux.
// Routes registered:
//   - GET /api/route       : compute bike routes between coordinates
//   - GET /api/all-ways    : fetch all ways from the database
//   - GET, POST /api/reviews : get or post reviews for ways
//   - POST /api/signup     : user signup
//   - POST /api/login      : user login
func RegisterRoutes(mux *http.ServeMux, handlers handler.Handlers) {
	mux.Handle("/api/route", corsMiddleware("GET", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlers.RouteHandler.HandleGetRoute()
		handler(w, r)
	})))

	mux.Handle("/api/all-ways", corsMiddleware("GET", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlers.WayHandler.HandleAllWays()
		handler(w, r)
	})))

	mux.Handle("/api/reviews", corsMiddleware("GET", "POST", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			handler := handlers.ReviewHandler.HandleGetReviews()
			handler(w, r)
		case http.MethodPost:
			handler := handlers.ReviewHandler.HandlePostReview()
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
			handler := handlers.AuthHandler.HandleSignup()
			handler(w, r)
		}
	})))

	mux.Handle("/api/login", corsMiddleware("POST", "OPTIONS")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			handler := handlers.AuthHandler.HandleLogin()
			handler(w, r)
		}
	})))
}
