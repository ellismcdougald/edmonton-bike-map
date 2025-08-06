package routing

import "github.com/ellismcdougald/edmonton-bike-map/pkg/model"

type Router interface {
	FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (dist float64, path []int64)
}
