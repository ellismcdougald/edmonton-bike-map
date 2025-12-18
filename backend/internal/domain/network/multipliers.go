package network

import (
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

const (
	tagBicycle      = "bicycle"
	tagBike         = "bike"
	tagCycleway     = "cycleway"
	tagLCN          = "lcn"
	tagMotorVehicle = "motor_vehicle"
	tagOneway       = "one_way"
	tagMTB          = "mtb:scale:imba"
)

// computeTagsMultiplier computes a multiplier for a way's weight based on its tags
func computeTagsMultiplier(tags map[string]string) float64 {
	highwayMultiplier := computeHighwayMultiplier(tags["highway"])
	surfaceMultiplier := computeSurfaceMultiplier(tags["surface"])
	bikeFriendlyMultiplier := computeBikeFriendlyMultiplier(tags)

	// Do not punish non-cycleways if they are cycle designated
	if bikeFriendlyMultiplier < 1 && highwayMultiplier > 1 {
		highwayMultiplier = 1
	}

	return highwayMultiplier * surfaceMultiplier * bikeFriendlyMultiplier
}

// computeHighwayMultiplier computes a multiplier based on the highway tag
func computeHighwayMultiplier(highwayTag string) float64 {
	highwayPenalty := map[string]float64{
		"cycleway":    0.9,
		"path":        0.9,
		"residential": 1,
		"tertiary":    1.2,
		"secondary":   1.5,
		"primary":     2.0,
	}

	highwayMultiplier, found := highwayPenalty[highwayTag]
	if !found {
		highwayMultiplier = 1.5
	}

	return highwayMultiplier
}

// computeSurfaceMultiplier computes a multiplier based on the surface tag
func computeSurfaceMultiplier(surfaceTag string) float64 {
	surfacePenalty := map[string]float64{
		"asphalt": 0.95,
		"gravel":  1.1,
		"dirt":    1.2,
	}

	surfaceMultiplier, found := surfacePenalty[surfaceTag]
	if !found {
		surfaceMultiplier = 1.0
	}

	return surfaceMultiplier
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
	if tags[tagBicycle] == "no" || tags[tagBike] == "no" {
		bikeFriendlyMultiplier *= 3
	}
	if tags[tagMTB] != "" {
		bikeFriendlyMultiplier *= 10
	}
	if tags["highway"] == "steps" {
		bikeFriendlyMultiplier *= 5
	}
	return bikeFriendlyMultiplier
}

func computeReviewMultiplier(reviews []models.Review) float64 {
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
