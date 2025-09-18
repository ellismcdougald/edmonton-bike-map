package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// HandleRouteByCoordinates handles HTTP requests to compute a bike route between start and end coordinates.
// Query parameters: startLatitude, startLongitude, endLatitude, endLongitude (float64).
// Responds with a GeoJSON LineString representing the line path.
// Returns HTTP 400 if parameters are missing or invalid.
func (h *RealHandlers) HandleRouteByCoordinates() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		getFloatParam := func(query url.Values, paramName string) (result float64, err error) {
			param := query.Get(paramName)
			result, err = strconv.ParseFloat(param, 64)
			if err != nil {
				log.Printf("Error extracting parameter %s from query %v: %v", paramName, query, err)
				http.Error(writer, "Invalid "+paramName+err.Error(), http.StatusBadRequest)
				return 0, err
			}
			return
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

		dist, pathIds := h.Router.FindRouteFromCoordinates(h.Network, startLatitude, startLongitude, endLatitude, endLongitude)

		var coordinates = [][2]float64{}
		for _, nodeID := range pathIds {
			var nodeLonLat = [2]float64{}
			nodeLonLat[0] = h.Network.Nodes[nodeID].Longitude
			nodeLonLat[1] = h.Network.Nodes[nodeID].Latitude
			coordinates = append(coordinates, nodeLonLat)
		}

		geojson := map[string]any{
			"type": "Feature",
			"geometry": map[string]any{
				"type":        "LineString",
				"coordinates": coordinates,
			},
			"properties": map[string]any{
				"distance_km": dist,
			},
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(geojson)
		if err != nil {
			log.Printf("handleRoute error encoding json: %v", err)
		}
	}
}
