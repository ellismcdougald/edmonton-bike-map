package server

import (
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/internal/handler"
	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
)

// RegisterRoutes registers HTTP handlers for all API endpoints on the given ServeMux.
// Routes registered:
//   - POST /api/login            : user login
//   - POST /api/signup           : user signup
//   - GET /api/route             : compute bike routes between coordinates
//   - GET /api/all-ways          : fetch all ways from the database
//   - GET /api/nearest-way       : find the nearest way to given coordinates
//   - GET /api/way               : get way details by ID
//   - GET /api/adjacent-ways     : get ways adjacent to a given way
//   - GET, POST /api/reviews     : get or post reviews for ways
//   - POST /api/change-password  : change user password (auth required)
//   - GET, POST /api/settings    : get or update user settings (auth required)
//
// RegisterRoutes registers the server's API endpoints on the provided ServeMux,
// applying CORS to all routes and authentication middleware to protected routes.
//
// It sets up public endpoints (e.g., /api/login, /api/signup), guest-accessible endpoints
// (e.g., /api/route, /api/all-ways) that accept optional authentication, and protected
// endpoints that require authentication (e.g., /api/change-password, /api/settings, POST /api/reviews).
//
// mux is the HTTP ServeMux to attach routes to.
// handlers supplies the concrete handler implementations for each endpoint.
func RegisterRoutes(mux *http.ServeMux, handlers handler.Handlers) {
	// wrap handler with CORS and optionally Auth middleware
	wrap := func(handler http.Handler, methods []string, authType string) http.Handler {
		h := handler
		switch authType {
		case "required":
			h = middleware.AuthMiddleware(h)
		case "optional":
			h = middleware.OptionalAuthMiddleware(h)
		}
		h = middleware.CorsMiddleware(methods...)(h)
		return h
	}

	// PUBLIC:
	mux.Handle("/api/login", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleLogin()),
		[]string{"POST", "OPTIONS"},
		"none",
	))
	mux.Handle("/api/signup", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleSignup()),
		[]string{"POST", "OPTIONS"},
		"none",
	))

	// GUEST-ACCESSIBLE (optional auth):
	mux.Handle("/api/route", wrap(
		http.HandlerFunc(handlers.RouteHandler.HandleGetRoute()),
		[]string{"GET", "OPTIONS"},
		"optional",
	))
	mux.Handle("/api/all-ways", wrap(
		http.HandlerFunc(handlers.WayHandler.HandleAllWays()),
		[]string{"GET", "OPTIONS"},
		"optional",
	))
	mux.Handle("/api/nearest-way", wrap(
		http.HandlerFunc(handlers.WayHandler.HandleNearestWay()),
		[]string{"GET", "OPTIONS"},
		"optional",
	))
	mux.Handle("/api/way", wrap(
		http.HandlerFunc(handlers.WayHandler.HandleGetWay()),
		[]string{"GET", "OPTIONS"},
		"optional",
	))
	mux.Handle("/api/adjacent-ways", wrap(
		http.HandlerFunc(handlers.WayHandler.HandleGetAdjacentWays()),
		[]string{"GET", "OPTIONS"},
		"optional",
	))

	// REVIEWS (GET is optional, POST requires auth):
	mux.Handle("/api/reviews", wrap(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				// GET reviews is guest-accessible
				handlers.ReviewHandler.HandleGetReviews()(w, r)
			case http.MethodPost:
				// POST reviews requires authentication
				authHandler := middleware.AuthMiddleware(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						handlers.ReviewHandler.HandlePostReview()(w, r)
					},
				))
				authHandler.ServeHTTP(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		[]string{"GET", "POST", "OPTIONS"},
		"none", // We handle auth manually per method above
	))

	// PROTECTED (require auth):
	mux.Handle("/api/change-password", wrap(
		http.HandlerFunc(handlers.AuthHandler.HandleChangePassword()),
		[]string{"POST", "OPTIONS"},
		"required",
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
		"required",
	))
}
