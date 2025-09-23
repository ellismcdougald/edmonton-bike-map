package routing

import "github.com/ellismcdougald/edmonton-bike-map/pkg/model"

// Router defines route-finding functionality.
type Router interface {
	// FindRouteFromCoordinates finds the shortest route between start and end coordinates.
	FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (dist float64, path []int64)

	// FindRoutesFromCoordinates finds up to k distinct routes between start and end coordinates.
	// Returns a slice of model.Path, each containing Nodes and Cost.
	FindRoutesFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64, k int) []model.Path
}