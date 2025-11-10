package routing

import (
	"math"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// nearestNode returns the ID of the closest node in the network to the given latitude-longitude pair.
func nearestNode(latitude float64, longitude float64, nodes map[int64]models.Node) (nodeID int64) {
	smallestDistance := math.Inf(1)
	var smallestNodeID int64 = -1
	for id, node := range nodes {
		distance := squaredEucDistance(latitude, longitude, node.Latitude, node.Longitude)
		if distance < smallestDistance {
			smallestDistance = distance
			smallestNodeID = id
		}
	}
	return smallestNodeID
}

// squaredEucDistance computes the squared Euclidean distance between two latitude-longitude pairs.
func squaredEucDistance(lat1, lon1, lat2, lon2 float64) float64 {
	distLat := lat2 - lat1
	distLon := lon2 - lon1
	return distLat*distLat + distLon*distLon
}

