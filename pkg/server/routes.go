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
		log.Print("reviews triggered")
		handlePostReview(writer, request, db)
	})
}
