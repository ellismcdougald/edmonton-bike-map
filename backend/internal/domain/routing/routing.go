package routing

import "github.com/ellismcdougald/edmonton-bike-map/internal/models"

// FindRouteFromCoordinates returns the lowest-cost route between two latitude-longitude pairs.
func FindRouteFromCoordinates(network *models.Network, startLatitude, startLongitude, endLatitude, endLongitude float64) (dist float64, path []int64) {
	startNodeID := nearestNode(startLatitude, startLongitude, network.Nodes)
	endNodeID := nearestNode(endLatitude, endLongitude, network.Nodes)
	return findRoute(network, startNodeID, endNodeID)
}