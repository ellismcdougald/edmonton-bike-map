package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
)

// HandleAllWays returns an HTTP handler that responds with all ways as a GeoJSON FeatureCollection.
// It fetches all nodes and ways from the database via the respective services,
// constructs LineString geometries from node coordinates,
// and includes way tags as properties in each Feature.
// Returns HTTP 500 if any database fetch or JSON encoding error occurs,
// or if any way references a missing node.
func (h *RealHandlers) HandleAllWays() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		allNodes, err := h.NodeService.GetAllNodes() // map nodeId -> node
		if err != nil {
			log.Printf("Could not get nodes from db: %v", err)
			http.Error(writer, "Could not get nodes from database", http.StatusInternalServerError)
			return
		}
		allWays, err := h.WayService.GetAllWays() // array
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
}