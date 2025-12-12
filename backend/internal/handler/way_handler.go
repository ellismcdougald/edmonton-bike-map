package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/geo"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

type WayHandler struct {
	NodeService *service.NodeService
	WayService  *service.WayService
}

func NewWayHandler(nodeService *service.NodeService, wayService *service.WayService) *WayHandler {
	return &WayHandler{
		NodeService: nodeService,
		WayService:  wayService,
	}
}

// HandleAllWays returns an HTTP handler that fetches all ways and their associated nodes,
// maps them to GeoJSON features, and responds with a FeatureCollection in JSON format.
// Responds with HTTP 500 Internal Server Error if fetching data or mapping fails.
func (h *WayHandler) HandleAllWays() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		allNodes, err := h.NodeService.GetAllNodes()
		if err != nil {
			log.Printf("Could not get nodes from db: %v", err)
			http.Error(writer, "Could not get nodes from database", http.StatusInternalServerError)
			return
		}

		allWays, err := h.WayService.GetAllWays()
		if err != nil {
			log.Printf("Could not get ways from db: %v", err)
			http.Error(writer, "Could not get ways from database", http.StatusInternalServerError)
			return
		}

		allFeatures, err := geo.MapWaysToFeatures(allWays, allNodes)
		if err != nil {
			log.Printf("Error mapping ways to features: %v", err)
			http.Error(writer, "Error mapping ways to features: "+err.Error(), http.StatusInternalServerError)
			return
		}

		featureCollection := models.FeatureCollection{
			Type:     "FeatureCollection",
			Features: allFeatures,
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(featureCollection); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}

// HandleNearestWay returns an HTTP handler that finds the nearest way to given coordinates.
// Expects query parameters 'lat' and 'lng' (as floats).
// Responds with a simple JSON object containing the way ID and tags.
// Responds with HTTP 400 Bad Request if coordinates are invalid or missing.
// Responds with HTTP 404 Not Found if no ways are found nearby.
// Responds with HTTP 500 Internal Server Error if the database query fails.
func (h *WayHandler) HandleNearestWay() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		latStr := request.URL.Query().Get("lat")
		lngStr := request.URL.Query().Get("lng")

		if latStr == "" || lngStr == "" {
			http.Error(writer, "Missing required parameters: lat and lng", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			http.Error(writer, "Invalid latitude value", http.StatusBadRequest)
			return
		}

		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			http.Error(writer, "Invalid longitude value", http.StatusBadRequest)
			return
		}

		way, err := h.WayService.GetNearestWay(lat, lng)
		if err != nil {
			if errors.Is(err, service.ErrWayNotFound) {
				http.Error(writer, "No ways found nearby", http.StatusNotFound)
				return
			}

			log.Printf("Could not get nearest way: %v", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":   way.ID,
			"tags": way.Tags,
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}

// HandleGetWay returns an HTTP handler that fetches a single way by ID.
// Expects query parameter 'id' (as int64).
// Responds with a simple JSON object containing the way ID and tags.
// Responds with HTTP 400 Bad Request if ID is invalid or missing.
// Responds with HTTP 404 Not Found if the way is not found.
// Responds with HTTP 500 Internal Server Error if the database query fails.
func (h *WayHandler) HandleGetWay() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idStr := request.URL.Query().Get("id")

		if idStr == "" {
			http.Error(writer, "Missing required parameter: id", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(writer, "Invalid id value", http.StatusBadRequest)
			return
		}

		way, err := h.WayService.GetWay(id)
		if err != nil {
			if errors.Is(err, service.ErrWayNotFound) {
				http.Error(writer, "Way not found", http.StatusNotFound)
				return
			}

			log.Printf("Could not get way: %v", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":   way.ID,
			"tags": way.Tags,
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}
