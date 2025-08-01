// Package server serves the API endpoints for Edmonton Bike Map.
package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// RegisterRoutes registers HTTP routes on the given ServeMux.
// Endpoints:
// - /api/route: handles routing requests given latitude-longitude coordinates
func RegisterRoutes(mux *http.ServeMux, network *model.Graph, db *sql.DB) {
	mux.HandleFunc("/api/route", func(writer http.ResponseWriter, request *http.Request) {
		handleRouteByCoordinates(writer, request, network)
	})
	mux.HandleFunc("/api/all-ways", func(writer http.ResponseWriter, request *http.Request) {
		handleAllWays(writer, request, db)
	})
	mux.HandleFunc("/api/reviews", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		switch request.Method {
		case http.MethodOptions:
			writer.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			handleGetReviews(writer, request, db)
		case http.MethodPost:
			handlePostReview(writer, request, db)
		default:
			writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/signup", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		log.Print("Hit")
		switch request.Method {
		case http.MethodOptions:
			writer.WriteHeader(http.StatusNoContent)
		default:
			handleSignUp(writer, request, db)
		}
	})

	mux.HandleFunc("/api/login", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		switch request.Method {
		case http.MethodOptions:
			writer.WriteHeader(http.StatusNoContent)
		default:
			handleLogin(writer, request, db)
		}
	})
}
