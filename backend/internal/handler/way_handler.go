package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/geo"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

type WayHandler struct {
	NodeService *service.NodeService
	WayService *service.WayService
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

		allFeatures , err := geo.MapWaysToFeatures(allWays, allNodes)
		if err != nil {
			log.Printf("Error mapping ways to features: %v", err)
			http.Error(writer, "Error mapping ways to features: " + err.Error(), http.StatusInternalServerError)
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