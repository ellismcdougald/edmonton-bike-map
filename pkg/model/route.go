package model

// Route represents a bike route with an ID, duration in seconds, and a geometry line consisting of coordinates.
type Route struct {
	Id           int          // Unique identifier for the route
	Duration     int          // Duration of the route in seconds
	GeometryLine []Coordinate // Sequence of coordinates representing the route path
}

// Coordinate represents a geographic point with latitude and longitude.
type Coordinate struct {
	Latitude  float64 // Latitude in decimal degrees
	Longitude float64 // Longitude in decimal degrees
}
