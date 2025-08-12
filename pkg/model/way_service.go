package model

// WayService defines the interface for managing ways in a data store.
// It supports inserting an indivdiual way and getting all ways.
type WayService interface {
	Insert(w DBWay) error
	GetAllWays() ([]DBWay, error)
}
