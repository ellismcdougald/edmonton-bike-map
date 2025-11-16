package handler

// Handlers groups all HTTP handlers.
type Handlers struct {
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
	WayHandler    *WayHandler
	ReviewHandler *ReviewHandler
	RouteHandler  *RouteHandler
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	wayHandler *WayHandler,
	reviewHandler *ReviewHandler,
	routeHandler *RouteHandler,
) *Handlers {
	return &Handlers{
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		WayHandler:    wayHandler,
		ReviewHandler: reviewHandler,
		RouteHandler:  routeHandler,
	}
}
