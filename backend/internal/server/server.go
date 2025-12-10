package server

import (
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/internal/handler"
	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
)

// RegisterRoutes registers HTTP handlers for all API endpoints on the given ServeMux.
// Routes registered:
//   - GET /api/route       : compute bike routes between coordinates
//   - GET /api/all-ways    : fetch all ways from the database
//   - GET, POST /api/reviews : get or post reviews for ways
//   - POST /api/signup     : user signup
// RegisterRoutes registers the server's API endpoints on the provided ServeMux,
// applying CORS to all routes and authentication middleware to protected routes.
//
// It sets up public endpoints (e.g., /api/login, /api/signup) and protected
// endpoints that require authentication (e.g., /api/route, /api/all-ways,
// /api/reviews, /api/change-password, /api/settings).
//
// mux is the HTTP ServeMux to attach routes to.
// handlers supplies the concrete handler implementations for each endpoint.
func RegisterRoutes(mux *http.ServeMux, handlers handler.Handlers) {
	// wrap handler with CORS and optionally Auth middleware
	wrap := func(handler http.Handler, methods []string, requireAuth bool) http.Handler {
		h := handler
		if requireAuth {
			h = middleware.AuthMiddleware(h)
		}
		h = middleware.CorsMiddleware(methods...)(h)
		return h
	}

	// PUBLIC:
	mux.Handle("/api/login", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleLogin()),
		[]string{"POST", "OPTIONS"},
		false,
	))
	mux.Handle("/api/signup", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleSignup()),
		[]string{"POST", "OPTIONS"},
		false,
	))

	// PROTECTED:
	mux.Handle("/api/route", wrap(
		http.HandlerFunc(handlers.RouteHandler.HandleGetRoute()),
		[]string{"GET", "OPTIONS"},
		true,
	))
	mux.Handle("/api/all-ways", wrap(
		http.HandlerFunc(handlers.WayHandler.HandleAllWays()),
		[]string{"GET", "OPTIONS"},
		true,
	))
	mux.Handle("/api/reviews", wrap(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.ReviewHandler.HandleGetReviews()(w, r)
			case http.MethodPost:
				handlers.ReviewHandler.HandlePostReview()(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		[]string{"GET", "POST", "OPTIONS"},
		true,
	))
	mux.Handle("/api/change-password", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleChangePassword()),
		[]string{"POST", "OPTIONS"},
		true,
	))
	mux.Handle("/api/settings", wrap(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.UserHandler.HandleGetSettings()(w, r)
			case http.MethodPost:
				handlers.UserHandler.HandleUpdateCyclingSpeed()(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		[]string{"GET", "POST", "OPTIONS"},
		true,
	))
}