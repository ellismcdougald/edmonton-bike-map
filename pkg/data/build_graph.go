// Package data handles fetching OSM data and parsing it into a graph structure suitable for bike routing
package data

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// OSMResponse consists of a list of OSMElements.
type OSMResponse struct {
	Elements []OSMElement `json:"elements"`
}

// OSMElement represents OSM 'nodes' and 'ways' with associated tags and coordinates.
type OSMElement struct {
	Type  string            `json:"type"` // "node" or "way"
	ID    int64             `json:"id"`
	Lat   float64           `json:"lat,omitempty"`   // only nodes
	Lon   float64           `json:"lon,omitempty"`   // only nodes
	Nodes []int64           `json:"nodes,omitempty"` // only ways
	Tags  map[string]string `json:"tags,omitempty"`
}

// FeatureCollection represents a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Feature represents a GeoJSON feature with associated properties and geometry.
type Feature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   Geometry               `json:"geometry"`
}

// Geometry represents the GeoJSON geometry of a feature.
type Geometry struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

// GetAllGeoJsonData transforms the data into GeoJSON format
func GetAllGeoJsonData(filename string) ([]byte, error) {
	resp, err := parseOSMJSON(filename)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[int64]OSMElement)
	for _, el := range resp.Elements {
		if el.Type == "node" {
			nodeMap[el.ID] = el
		}
	}

	var allFeatures []Feature
	for _, el := range resp.Elements {
		if el.Type == "way" {
			var coordinates = [][2]float64{}
			for _, nodeID := range el.Nodes {
				var nodeLonLat = [2]float64{}
				nodeLonLat[0] = nodeMap[nodeID].Lon
				nodeLonLat[1] = nodeMap[nodeID].Lat
				coordinates = append(coordinates, nodeLonLat)
			}

			wayFeature := Feature{
				Type: "Feature",
				Properties: map[string]any{
					"name": el.Tags["name"],
				},
				Geometry: Geometry{
					Type:        "LineString",
					Coordinates: coordinates,
				},
			}
			allFeatures = append(allFeatures, wayFeature)
		}
	}

	featureCollection := FeatureCollection{
		Type:     "FeatureCollection",
		Features: allFeatures,
	}

	jsonData, err := json.MarshalIndent(featureCollection, "", "  ")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("web/edmonton_bike_data_geo.json", jsonData, 0644)
	if err != nil {
		log.Printf("Error")
		return nil, err
	} else {
		log.Printf("Success")
	}

	return jsonData, nil
}

// BuildGraph reads OSM json data from a file and parses it into a graph structure
func BuildGraph(filename string) (*model.Graph, error) {
	resp, err := parseOSMJSON(filename)
	if err != nil {
		return nil, err
	}

	network := model.Graph{
		Nodes: make(map[int64]model.Node),
		Edges: make(map[int64][]model.Edge),
	}
	for _, el := range resp.Elements {
		if el.Type == "node" {
			network.Nodes[el.ID] = model.Node{
				Latitude:  el.Lat,
				Longitude: el.Lon,
			}
		}
	}
	for j, el := range resp.Elements {
		if j == 0 {
			fmt.Printf("%+v", el)
		}
		if el.Type == "way" {
			for i := 0; i < len(el.Nodes)-1; i++ {
				fromID := el.Nodes[i]
				fromCoord := network.Nodes[fromID]
				toID := el.Nodes[i+1]
				toCoord := network.Nodes[toID]

				dist := haversineDistance(fromCoord.Latitude, fromCoord.Longitude, toCoord.Latitude, toCoord.Longitude)
				weight := computeWayWeight(dist, el.Tags)

				network.Edges[fromID] = append(network.Edges[fromID], model.Edge{
					To:     toID,
					Weight: weight,
				})
				network.Edges[toID] = append(network.Edges[toID], model.Edge{
					To:     fromID,
					Weight: weight,
				})
			}
		}
	}
	return &network, nil
}

// computeWayWeight computes a weight for an edge by modifying the distance using the tags
func computeWayWeight(distance float64, tags map[string]string) float64 {
	highwayPenalty := map[string]float64{
		"cycleway":    0.5,
		"residential": 1,
		"tertiary":    1.2,
		"secondary":   1.5,
		"primary":     2.0,
		"motorway":    math.Inf(1),
		"trunk":       math.Inf(1),
	}
	surfacePenalty := map[string]float64{
		"concrete": 1.0,
		"asphalt":  1.0,
		"gravel":   1.5,
		"dirt":     1.75,
	}

	highwayMultiplier, found := highwayPenalty[tags["highway"]]
	if !found {
		highwayMultiplier = 1.5
	}

	surfaceMultiplier, found := surfacePenalty[tags["surface"]]
	if !found {
		surfaceMultiplier = 1.75
	}

	bikeFriendlyMultiplier := 1.0
	if tags["cycleway"] != "" {
		bikeFriendlyMultiplier *= 0.8
	}
	if tags["bicycle"] == "designated" {
		bikeFriendlyMultiplier *= 0.9
	} else if tags["bicycle"] == "yes" || tags["bike"] == "yes" {
		bikeFriendlyMultiplier *= 0.9
	}
	if tags["motor_vehicle"] == "no" {
		bikeFriendlyMultiplier *= 0.7
	}
	if tags["lcn"] == "yes" {
		bikeFriendlyMultiplier *= 0.9
	}

	return distance * highwayMultiplier * surfaceMultiplier * bikeFriendlyMultiplier
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // meters

	degToRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	lat1Rad := degToRad(lat1)
	lat2Rad := degToRad(lat2)
	deltaLat := degToRad(lat2 - lat1)
	deltaLon := degToRad(lon2 - lon1)

	sinDeltaLat := math.Sin(deltaLat / 2)
	sinDeltaLon := math.Sin(deltaLon / 2)

	a := sinDeltaLat*sinDeltaLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinDeltaLon*sinDeltaLon
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

func parseOSMJSON(filename string) (*OSMResponse, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var resp OSMResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
