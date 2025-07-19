// Package server serves the API endpoints for Edmonton Bike Map.
package server

import (
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// RegisterRoutes registers HTTP routes on the given ServeMux.
// Endpoints:
// - /api/route: handles routing requests given latitude-longitude coordinates
func RegisterRoutes(mux *http.ServeMux, network *model.Graph) {
	mux.HandleFunc("/api/route", func(writer http.ResponseWriter, request *http.Request) {
		handleRouteByCoordinates(writer, request, network)
	})
}
