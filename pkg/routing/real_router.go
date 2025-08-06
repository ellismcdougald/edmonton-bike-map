package routing

import "github.com/ellismcdougald/edmonton-bike-map/pkg/model"

type RealRouter struct{}

func (r RealRouter) FindRouteFromCoordinates(network *model.Graph, startLatitude, startLongitude, endLatitude, endLongitude float64) (float64, []int64) {
	return FindRouteFromCoordinates(network, startLatitude, startLongitude, endLatitude, endLongitude)
}
