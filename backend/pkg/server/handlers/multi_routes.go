package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// HandleRoutesByCoordinates handles HTTP requests to compute up to k bike routes between start and end coordinates.
// Query parameters: startLatitude, startLongitude, endLatitude, endLongitude (float64), k (optional, defaults to 3).
// Responds with a GeoJSON FeatureCollection containing LineStrings for each route.
// Returns HTTP 400 if parameters are missing or invalid.
func (h *RealHandlers) HandleRoutesByCoordinates() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		getFloatParam := func(query url.Values, paramName string) (float64, error) {
			param := query.Get(paramName)
			result, err := strconv.ParseFloat(param, 64)
			if err != nil {
				log.Printf("Error extracting parameter %s from query %v: %v", paramName, query, err)
				http.Error(writer, "Invalid "+paramName+": "+err.Error(), http.StatusBadRequest)
				return 0, err
			}
			return result, nil
		}

		startLatitude, err := getFloatParam(query, "startLatitude")
		if err != nil {
			return
		}
		startLongitude, err := getFloatParam(query, "startLongitude")
		if err != nil {
			return
		}
		endLatitude, err := getFloatParam(query, "endLatitude")
		if err != nil {
			return
		}
		endLongitude, err := getFloatParam(query, "endLongitude")
		if err != nil {
			return
		}

		// Optional k parameter
		k := 3
		if kStr := query.Get("k"); kStr != "" {
			if kParsed, err := strconv.Atoi(kStr); err == nil && kParsed > 0 {
				k = kParsed
			}
		}

		routes := h.Router.FindRoutesFromCoordinates(h.Network, startLatitude, startLongitude, endLatitude, endLongitude, k)
		if len(routes) == 0 {
			http.Error(writer, "No routes found", http.StatusNotFound)
			return
		}

		var features []any
		for _, route := range routes {
			var coordinates = [][2]float64{}
			for _, nodeID := range route.Nodes {
				node := h.Network.Nodes[nodeID]
				coordinates = append(coordinates, [2]float64{node.Longitude, node.Latitude})
			}

			feature := map[string]any{
				"type": "Feature",
				"geometry": map[string]any{
					"type":        "LineString",
					"coordinates": coordinates,
				},
				"properties": map[string]any{
					"distance_km": route.Cost,
				},
			}
			features = append(features, feature)
		}

		geojson := map[string]any{
			"type":     "FeatureCollection",
			"features": features,
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(geojson); err != nil {
			log.Printf("HandleRoutesByCoordinates error encoding json: %v", err)
		}
	}
}
