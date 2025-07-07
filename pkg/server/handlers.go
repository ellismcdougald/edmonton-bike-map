package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/routing"
)

func handleRoute(writer http.ResponseWriter, request *http.Request, network *model.Graph) {
	query := request.URL.Query()

	startStr := query.Get("start")
	endStr := query.Get("end")

	startID, _ := strconv.ParseInt(startStr, 10, 64)
	endID, _ := strconv.ParseInt(endStr, 10, 64)
	_, pathIds := routing.FindRoute(network, startID, endID)

	var coordinates = [][2]float64{}
	for _, nodeID := range pathIds {
		var nodeLonLat = [2]float64{}
		nodeLonLat[0] = network.Nodes[nodeID].Longitude
		nodeLonLat[1] = network.Nodes[nodeID].Latitude
		coordinates = append(coordinates, nodeLonLat)
	}

	geojson := map[string]any {
		"type": "Feature",
		"geometry": map[string]any {
			"type": "LineString",
			"coordinates": coordinates,
		},
		"properties": map[string]any{},
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(geojson)
}