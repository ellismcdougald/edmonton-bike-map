package geo

import (
	"fmt"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

func MapWaysToFeatures(allWays []models.Way, allNodes map[int64]models.Node) ([]models.Feature, error) {
	var features []models.Feature

    for _, way := range allWays {
        coords := make([][2]float64, 0, len(way.NodeIDs))
        for _, nodeID := range way.NodeIDs {
            node, ok := allNodes[nodeID]
            if !ok {
                return nil, fmt.Errorf("node %d missing from node data", nodeID)
            }
            coords = append(coords, [2]float64{node.Longitude, node.Latitude})
        }

        geometry := models.Geometry{
          Type:        "LineString",
          Coordinates: coords,
        }

        props := make(map[string]string, len(way.Tags)+1)
        for k, v := range way.Tags {
            props[k] = v
        }
        props["id"] = strconv.FormatInt(way.ID, 10)

        feature := models.Feature{
            Type:       "Feature",
            Properties: props,
            Geometry:   geometry,
        }

        features = append(features, feature)
    }

    return features, nil
	}