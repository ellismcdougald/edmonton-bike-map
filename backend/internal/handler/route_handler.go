package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
)

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
		// Decode user id
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(writer, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("user_handler: token validation failed: %v", err)
			http.Error(writer, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

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

		// call service
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

		// build coordinates as [][2]float64 with [lon, lat] like the old handler
		var coordinates = [][2]float64{}
		for _, n := range nodes {
			var lonLat [2]float64
			lonLat[0] = n.Longitude
			lonLat[1] = n.Latitude
			coordinates = append(coordinates, lonLat)
		}

		cyclingSpeed, err := h.UserService.GetCyclingSpeed(claims.UserID)
		if err != nil {
			log.Printf("error getting cycling speed: %v", err)
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			return
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
