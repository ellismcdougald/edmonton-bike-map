package routing

import "github.com/ellismcdougald/edmonton-bike-map/pkg/model"

// RealRouter implements the Router interface.
// It provides route finding using FindRouteFromCoordinates.
type RealRouter struct{}

// FindRouteFromCoordinates implements the Router interface method for RealRouter.
// It finds the shortest route between start and end coordinates using the given network.
func (r RealRouter) FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (float64, []int64) {
	return FindRouteFromCoordinates(network, startLatitude, startLongitude, endLatitude, endLongitude)
}
