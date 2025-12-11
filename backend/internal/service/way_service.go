package service

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// WayService provides operations related to ways.
type WayService struct {
	WayRepository repository.WayRepository
}

// NewWayService creates a new instance of WayService.
func NewWayService(wayRepo repository.WayRepository) *WayService {
	return &WayService{
		WayRepository: wayRepo,
	}
}

// InsertWay inserts a new way using the WayRepository.
func (s *WayService) InsertWay(way models.Way) error {
	return s.WayRepository.Insert(way)
}

// GetAllWays retrieves all ways using the WayRepository.
func (s *WayService) GetAllWays() ([]models.Way, error) {
	return s.WayRepository.GetAllWays()
}

// GetNearestWay retrieves the nearest way to the given coordinates using the WayRepository.
func (s *WayService) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	return s.WayRepository.GetNearestWay(latitude, longitude)
}
