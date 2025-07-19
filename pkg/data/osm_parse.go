package data

import (
	"encoding/json"
	"log"
	"os"
)

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

// parseOSMJSON parses OSM data from a file into the OSMResponse struct
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
