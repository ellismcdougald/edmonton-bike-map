// Package data handles fetching OSM data and parsing it into a graph structure suitable for bike routing
package data

import (
	"log"
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

const (
	tagBicycle      = "bicycle"
	tagBike         = "bike"
	tagCycleway     = "cycleway"
	tagHighway      = "highway"
	tagLCN          = "lcn"
	tagMotorVehicle = "motor_vehicle"
	tagOneway       = "one_way"
)

// BuildGraph reads OSM json data from a file and parses it into a graph structure
func BuildGraph(filename string) (*model.Graph, error) {
	resp, err := ParseOSMJSON(filename)
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
	for _, el := range resp.Elements {
		if el.Type == "way" {
			for i := 0; i < len(el.Nodes)-1; i++ {
				fromID := el.Nodes[i]
				fromCoord, fromExists := network.Nodes[fromID]
				toID := el.Nodes[i+1]
				toCoord, toExists := network.Nodes[toID]
				if !fromExists || !toExists {
					log.Printf("Warning: node missing for way %d: %d to %d", el.ID, fromID, toID)
					continue
				}

				dist := haversineDistance(fromCoord.Latitude, fromCoord.Longitude, toCoord.Latitude, toCoord.Longitude)
				weight := computeWayWeight(dist, el.Tags)

				network.Edges[fromID] = append(network.Edges[fromID], model.Edge{
					To:     toID,
					Weight: weight,
				})
				if el.Tags[tagOneway] != "yes" {
					network.Edges[toID] = append(network.Edges[toID], model.Edge{
						To:     fromID,
						Weight: weight,
					})
				}
			}
		}
	}
	return &network, nil
}

// computeWayWeight computes a weight for an edge by modifying the distance using the tags
func computeWayWeight(distance float64, tags map[string]string) float64 {
	highwayPenalty := map[string]float64{
		"cycleway":    0.9,
		"residential": 1,
		"tertiary":    1.2,
		"secondary":   1.5,
		"primary":     2.0,
		"motorway":    math.Inf(1),
		"trunk":       math.Inf(1),
	}
	/*
		TODO: decide whether to proceed with surface penalties. Not all data is marked with surface
		// Would need to infer "asphalt" for good bike connections where they are not labeled

		surfacePenalty := map[string]float64{
			"concrete": 1.0,
			"asphalt":  1.0,
			"gravel":   1.5,
			"dirt":     1.75,
		}

		surfaceMultiplier, found := surfacePenalty[tags["surface"]]
		if !found {
			surfaceMultiplier = 1.75
		}
	*/

	highwayMultiplier, found := highwayPenalty[tags[tagHighway]]
	if !found {
		highwayMultiplier = 1.5
	}

	bikeFriendlyMultiplier := computeBikeFriendlyMultiplier(tags)

	// Do not punish non-cycleways if they are cycle designated
	if bikeFriendlyMultiplier < 1 && highwayMultiplier > 1 {
		highwayMultiplier = 1
	}

	return distance * highwayMultiplier * bikeFriendlyMultiplier
}

// computeBikeFriendlyMultiplier computes a multiplier for a way's weight based on the bike characteristics in its tags
func computeBikeFriendlyMultiplier(tags map[string]string) float64 {
	bikeFriendlyMultiplier := 1.0
	if tags[tagCycleway] != "" || tags[tagBicycle] == "designated" || tags[tagMotorVehicle] == "no" {
		bikeFriendlyMultiplier *= 0.9
	}
	if tags[tagBicycle] == "yes" || tags[tagBike] == "yes" || tags[tagLCN] == "yes" {
		bikeFriendlyMultiplier *= 0.95
	}
	return bikeFriendlyMultiplier
}

// haversineDistance computes the haversine distance between two latitude-longitude coordinate pairs
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

	return earthRadius * c
}
