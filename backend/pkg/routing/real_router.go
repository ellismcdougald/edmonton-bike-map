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

// FindRoutesFromCoordinates finds up to k distinct routes between start and end coordinates.
// Returns a slice of model.Path, each containing Nodes and Cost.
func (r RealRouter) FindRoutesFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64, k int) []model.Path {
	return FindRoutesFromCoordinates(network, startLatitude, startLongitude, endLatitude, endLongitude, k)
}
