package model

type WayService interface {
	Insert(w DBWay) error
	GetAllWays() ([]DBWay, error)
}
