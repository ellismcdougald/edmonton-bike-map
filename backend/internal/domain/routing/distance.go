package routing

import (
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// estimateDistance returns the estimated distance of the given path in kilometres
func estimateDistance(path []int64, nodes map[int64]models.Node) float64 {
	if len(path) < 2 {
		return 0
	}

	dist := 0.0
	for i := 1; i < len(path); i++ {
		node1, ok1 := nodes[path[i-1]]
		node2, ok2 := nodes[path[i]]
		if !ok1 || !ok2 {
			panic("estimateDistance: node ID not found in graph")
		}
		segDist := planarDistance(node1.Latitude, node1.Longitude, node2.Latitude, node2.Longitude)
		dist += segDist
	}

	return dist
}

// planarDistance returns distance in km
func planarDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const metersPerDeg = 111320.0 // approximate meters per degree
	latAvg := (lat1 + lat2) / 2 * math.Pi / 180

	dx := (lon2 - lon1) * math.Cos(latAvg) * metersPerDeg
	dy := (lat2 - lat1) * metersPerDeg

	return math.Sqrt(dx*dx+dy*dy) / 1000.0 // convert to km
}