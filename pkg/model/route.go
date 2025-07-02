package model

type Route struct {
	Id int
	Duration int
	GeometryLine []Coordinate
}

type Coordinate struct {
	Latitude float64
	Longitude float64
}