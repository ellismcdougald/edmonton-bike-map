package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/routing"
)

// APIHandlers defines the set of HTTP handler functions for the API endpoints,
// including user authentication, routing, ways, and reviews.
type APIHandlers interface {
	HandleLogin() http.HandlerFunc
	HandleSignup() http.HandlerFunc
	HandleRouteByCoordinates() http.HandlerFunc
	HandleAllWays() http.HandlerFunc
	HandleGetReviews() http.HandlerFunc
	HandlePostReview() http.HandlerFunc
}

// RealHandlers implements the APIHandlers interface, providing concrete
// HTTP handlers wired to the underlying services and data sources.
type RealHandlers struct {
	UserService   model.UserService
	NodeService   model.NodeService
	WayService    model.WayService
	ReviewService model.ReviewService
	DB            *sql.DB
	Network       *model.Graph
	Router        routing.Router
}
