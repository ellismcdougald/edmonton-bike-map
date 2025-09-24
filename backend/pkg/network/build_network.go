package network

import (
	"fmt"
	"log"
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

func BuildNetwork(nodeService model.NodeService, wayService model.WayService, reviewService model.ReviewService) (*model.Graph, error) {
	allDBNodes, err := nodeService.GetAllNodes()
	if err != nil {
		return nil, fmt.Errorf("error getting all nodes: %w", err)
	}
	allNodes := make(map[int64]model.Node, len(allDBNodes))
	for _, node := range allDBNodes {
		allNodes[node.ID] = model.Node{Latitude: node.Latitude, Longitude: node.Longitude}
	}

	allWays, err := wayService.GetAllWays()
	if err != nil {
		return nil, fmt.Errorf("error getting all ways: %w", err)
	}

	allReviews, err := reviewService.GetAllReviews()
	if err != nil {
		return nil, fmt.Errorf("error getting all reviews: %w", err)
	}

	edgesByNodeId := make(map[int64][]model.Edge, len(allNodes))
	for _, way := range allWays {
		tagsMultiplier := computeTagsMultiplier(way.Tags)
		reviewMultiplier := 1.0
		wayReviews, ok := allReviews[int(way.ID)]
		if ok {
			reviewMultiplier = computeReviewMultiplier(wayReviews)
		}

		for i := 0; i < len(way.NodeIDs)-1; i++ {
			fromID := way.NodeIDs[i]
			fromCoord, fromExists := allNodes[fromID]
			toID := way.NodeIDs[i+1]
			toCoord, toExists := allNodes[toID]
			if !fromExists || !toExists {
				log.Printf("Warning: node missing for way %d: %d to %d", way.ID, fromID, toID)
				continue
			}

			dist := haversineDistance(fromCoord.Latitude, fromCoord.Longitude, toCoord.Latitude, toCoord.Longitude)
			weight := dist * tagsMultiplier * reviewMultiplier

			edgesByNodeId[fromID] = append(edgesByNodeId[fromID], model.Edge{
				To:     toID,
				Weight: weight,
			})
			if way.Tags[tagOneway] != "yes" {
				edgesByNodeId[toID] = append(edgesByNodeId[toID], model.Edge{
					To:     fromID,
					Weight: weight,
				})
			}
		}
	}

	network := model.Graph{
		Nodes: allNodes,
		Edges: edgesByNodeId,
	}
	return &network, nil
}

func computeReviewMultiplier(reviews []model.Review) float64 {
	if len(reviews) == 0 {
		return 1.0
	}

	total := 0
	for _, review := range reviews {
		total += review.Rating
	}
	average := float64(total) / float64(len(reviews))

	multiplier := 1.2 - 0.1*average

	// Consider number of reviews in strength of multiplier
	confidence := math.Min(1.0, float64(len(reviews))/10.0)

	return 1.0 + (multiplier-1.0)*confidence
}

const (
	tagBicycle      = "bicycle"
	tagBike         = "bike"
	tagCycleway     = "cycleway"
	tagLCN          = "lcn"
	tagMotorVehicle = "motor_vehicle"
	tagOneway       = "one_way"
)

func computeTagsMultiplier(tags map[string]string) float64 {
	highwayMultiplier := computeHighwayMultiplier(tags["highway"])
	bikeFriendlyMultiplier := computeBikeFriendlyMultiplier(tags)

	// Do not punish non-cycleways if they are cycle designated
	if bikeFriendlyMultiplier < 1 && highwayMultiplier > 1 {
		highwayMultiplier = 1
	}

	return highwayMultiplier * bikeFriendlyMultiplier
}

func computeHighwayMultiplier(highwayTag string) float64 {
	highwayPenalty := map[string]float64{
		"cycleway":    0.9,
		"path": 0.9,
		"residential": 1,
		"tertiary":    1.2,
		"secondary":   1.5,
		"primary":     2.0,
		"motorway":    math.Inf(1),
		"trunk":       math.Inf(1),
	}

	highwayMultiplier, found := highwayPenalty[highwayTag]
	if !found {
		highwayMultiplier = 1.5
	}

	return highwayMultiplier
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

// TODO: decide whether to proceed with surface penalties. Not all data is marked with surface
// Would need to infer "asphalt" for good bike connections where they are not labeled
// func computeSurfaceMultiplier(surfaceTag string) float64 {
//

// 	surfacePenalty := map[string]float64{
// 		"concrete": 1.0,
// 		"asphalt":  1.0,
// 		"gravel":   1.5,
// 		"dirt":     1.75,
// 	}

// 	surfaceMultiplier, found := surfacePenalty[surfaceTag]
// 	if !found {
// 		surfaceMultiplier = 1.75
// 	}

// 	return surfaceMultiplier
// }
