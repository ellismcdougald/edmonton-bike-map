package service

import (
	"errors"

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

// GetWay retrieves a way by ID using the WayRepository.
func (s *WayService) GetWay(id int64) (*models.Way, error) {
	way, err := s.WayRepository.GetWay(id)
	if err != nil {
		if errors.Is(err, repository.ErrWayNotFound) {
			return nil, ErrWayNotFound
		}
		return nil, err
	}
	return way, nil
}

// GetNearestWay retrieves the nearest way to the given coordinates using the WayRepository.
func (s *WayService) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	way, err := s.WayRepository.GetNearestWay(latitude, longitude)
	if err != nil {
		if errors.Is(err, repository.ErrWayNotFound) {
			return nil, ErrWayNotFound
		}
		return nil, err
	}
	return way, nil
}
