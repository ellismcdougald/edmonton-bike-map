// Package data handles fetching OSM data and parsing it into a graph structure
package data

import (
	"encoding/json"
	"math"
	"os"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// OSMResponse consists of a list of OSMElements
type OSMResponse struct {
    Elements []OSMElement `json:"elements"`
}

// OSMElement represents OSM 'nodes' and 'ways'
type OSMElement struct {
    Type  string            `json:"type"`            // "node" or "way"
    ID    int64             `json:"id"`
    Lat   float64           `json:"lat,omitempty"`   // only nodes
    Lon   float64           `json:"lon,omitempty"`   // only nodes
    Nodes []int64           `json:"nodes,omitempty"` // only ways
    Tags  map[string]string `json:"tags,omitempty"`
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
				Latitude: el.Lat,
				Longitude: el.Lon,
			}
		}
	}
	for _, el := range resp.Elements {
		if el.Type == "way" {
			for i := 0; i < len(el.Nodes) - 1; i++ {
				fromID := el.Nodes[i]
				fromCoord := network.Nodes[fromID]
				toID := el.Nodes[i + 1]
				toCoord := network.Nodes[toID]
				
				dist := haversineDistance(fromCoord.Latitude, fromCoord.Longitude, toCoord.Latitude, toCoord.Longitude)

				network.Edges[fromID] = append(network.Edges[fromID], model.Edge{
					To: toID,
					Weight: dist,
				})
				network.Edges[toID] = append(network.Edges[toID], model.Edge{
					To: fromID,
					Weight: dist,
				})
			}
		}
	}
	return &network, nil
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


