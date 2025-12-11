package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

// DefaultCyclingSpeed is used for guest users
const DefaultCyclingSpeed = 15

type RouteHandler struct {
	RouteService *service.RouteService
	UserService  *service.UserService
}

func NewRouteHandler(routeService *service.RouteService, userService *service.UserService) *RouteHandler {
	return &RouteHandler{
		RouteService: routeService,
		UserService:  userService,
	}
}

func (h *RouteHandler) HandleGetRoute() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		getFloatParam := func(query url.Values, paramName string) (result float64, err error) {
			param := query.Get(paramName)
			result, err = strconv.ParseFloat(param, 64)
			if err != nil {
				log.Printf("Error extracting parameter %s from query %v: %v", paramName, query, err)
				http.Error(writer, "Invalid "+paramName+err.Error(), http.StatusBadRequest)
				return 0, err
			}
			return
		}

		startLatitude, err := getFloatParam(query, "startLatitude")
		if err != nil {
			return
		}
		startLongitude, err := getFloatParam(query, "startLongitude")
		if err != nil {
			return
		}
		endLatitude, err := getFloatParam(query, "endLatitude")
		if err != nil {
			return
		}
		endLongitude, err := getFloatParam(query, "endLongitude")
		if err != nil {
			return
		}

		dist, nodes, err := h.RouteService.FindRoute(startLatitude, startLongitude, endLatitude, endLongitude)
		if err != nil {
			log.Printf("error finding route: %v", err)
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			return
		}

		if dist < 0 || nodes == nil {
			http.Error(writer, "route not found", http.StatusNotFound)
			return
		}

		var coordinates = [][2]float64{}
		for _, n := range nodes {
			var lonLat [2]float64
			lonLat[0] = n.Longitude
			lonLat[1] = n.Latitude
			coordinates = append(coordinates, lonLat)
		}

		// Get cycling speed: use user's preferred speed if authenticated, otherwise use default
		cyclingSpeed := DefaultCyclingSpeed
		userID, ok := middleware.UserIDFromContext(request.Context())
		if ok {
			userSpeed, err := h.UserService.GetCyclingSpeed(userID)
			if err != nil {
				log.Printf("error getting cycling speed: %v", err)
				http.Error(writer, "internal server error", http.StatusInternalServerError)
				return
			}
			cyclingSpeed = userSpeed
		}
		timeMinutes := dist / float64(cyclingSpeed) * 60.0

		geojson := map[string]any{
			"type": "Feature",
			"geometry": map[string]any{
				"type":        "LineString",
				"coordinates": coordinates,
			},
			"properties": map[string]any{
				"distance_km":  dist,
				"time_minutes": timeMinutes,
			},
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(geojson); err != nil {
			log.Printf("handleRoute error encoding json: %v", err)
		}
	}
}
