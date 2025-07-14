package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/routing"
)

func handleRouteByCoordinates(writer http.ResponseWriter, request *http.Request, network *model.Graph) {
	query := request.URL.Query()

	getFloatParam := func(query url.Values, paramName string) (result float64, err error) {
		param := query.Get(paramName)
		result, err = strconv.ParseFloat(param, 64)
		if err != nil {
			log.Printf("Error extracting parameter %s from query %v: %v", paramName, query, err)
			http.Error(writer, "Invalid " + paramName + err.Error(), http.StatusBadRequest)
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

	// do something with the coordinates	
	_, pathIds := routing.FindRouteFromCoordinates(network, startLatitude, startLongitude, endLatitude, endLongitude)

	var coordinates = [][2]float64{}
	for _, nodeID := range pathIds {
		var nodeLonLat = [2]float64{}
		nodeLonLat[0] = network.Nodes[nodeID].Longitude
		nodeLonLat[1] = network.Nodes[nodeID].Latitude
		coordinates = append(coordinates, nodeLonLat)
	}

	geojson := map[string]any{
		"type": "Feature",
		"geometry": map[string]any{
			"type":        "LineString",
			"coordinates": coordinates,
		},
		"properties": map[string]any{},
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(geojson)
	if err != nil {
		log.Printf("handleRoute error encoding json: %v", err)
	}
}

/*
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

	geojson := map[string]any{
		"type": "Feature",
		"geometry": map[string]any{
			"type":        "LineString",
			"coordinates": coordinates,
		},
		"properties": map[string]any{},
	}

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(geojson)
	if err != nil {
		log.Printf("handleRoute error encoding json: %v", err)
	}
}
*/