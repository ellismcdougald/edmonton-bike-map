package routing

import "github.com/ellismcdougald/edmonton-bike-map/pkg/model"

// Router defines route-finding functionality.
// It finds the shortest route between start and end coordinates in a network.
type Router interface {
	FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (dist float64, path []int64)
}
