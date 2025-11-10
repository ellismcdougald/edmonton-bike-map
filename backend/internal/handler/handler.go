package handler

// Handlers groups all HTTP handlers.
type Handlers struct {
	AuthHandler   *AuthHandler
	WayHandler    *WayHandler
	ReviewHandler *ReviewHandler
	RouteHandler  *RouteHandler
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	authHandler *AuthHandler,
	wayHandler *WayHandler,
	reviewHandler *ReviewHandler,
	routeHandler *RouteHandler,
) *Handlers {
	return &Handlers{
		AuthHandler:   authHandler,
		WayHandler:    wayHandler,
		ReviewHandler: reviewHandler,
		RouteHandler:  routeHandler,
	}
}