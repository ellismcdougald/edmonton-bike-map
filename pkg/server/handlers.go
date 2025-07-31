package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/lib/pq" // <-- Postgres driver import

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/routing"
)

// handleRouteByCoordinates handles HTTP requests to compute a bike route between start and end coordinates.
// Query parameters: startLatitude, startLongitude, endLatitude, endLongitude (float64).
// Responds with a GeoJSON LineString representing the line path.
// Returns HTTP 400 if parameters are missing or invalid.
func handleRouteByCoordinates(writer http.ResponseWriter, request *http.Request, network *model.Graph) {
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

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

func handleAllWays(writer http.ResponseWriter, _ *http.Request, db *sql.DB) {
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	allNodes, err := model.GetAllNodes(db) // map nodeId -> node
	if err != nil {
		log.Printf("Could not get nodes from db: %v", err)
		http.Error(writer, "Could not get nodes from database", http.StatusInternalServerError)
		return
	}
	allWays, err := model.GetAllWays(db) // array
	if err != nil {
		log.Printf("Could not get ways from db: %v", err)
		http.Error(writer, "Could not get ways from database", http.StatusInternalServerError)
		return
	}

	allFeatures := []data.Feature{}

	for _, way := range allWays {
		coordinates := [][2]float64{}

		for _, nodeID := range way.NodeIDs {
			node, found := allNodes[nodeID]
			if !found {
				log.Printf("Error: node %d in way %d is missing from nodes data", nodeID, way.ID)
				http.Error(writer, "Node in way is missing from node data", http.StatusInternalServerError)
				return
			}

			coordinates = append(coordinates, [2]float64{node.Longitude, node.Latitude})
		}

		geometry := data.Geometry{
			Type:        "LineString",
			Coordinates: coordinates,
		}

		way.Tags["id"] = strconv.FormatInt(way.ID, 10)

		feature := data.Feature{
			Type:       "Feature",
			Properties: way.Tags,
			Geometry:   geometry,
		}

		allFeatures = append(allFeatures, feature)
	}

	featureCollection := data.FeatureCollection{
		Type:     "FeatureCollection",
		Features: allFeatures,
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(featureCollection)
	if err != nil {
		log.Printf("Error encoding json")
		http.Error(writer, "Could not encode json", http.StatusInternalServerError)
	}
}

func handlePostReview(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var review model.Review
	if err := json.NewDecoder(request.Body).Decode(&review); err != nil {
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	if review.Rating < 1 || review.Rating > 10 || review.ReviewText == "" {
		http.Error(writer, "Invalid review: Rating must be between 1 and 10 inclusive", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
	INSERT INTO reviews (way_id, rating, comment)
	VALUES ($1, $2, $3)
`, review.WayID, review.Rating, review.ReviewText)
	if err != nil {
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}

}
