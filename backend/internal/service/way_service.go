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

// GetAdjacentWays retrieves all ways that share at least one node with the given way.
// Excludes the original way from the results.
func (s *WayService) GetAdjacentWays(wayID int64) ([]models.Way, error) {
	// Get the original way to obtain its node IDs
	way, err := s.GetWay(wayID)
	if err != nil {
		return nil, err
	}

	if len(way.NodeIDs) == 0 {
		return []models.Way{}, nil
	}

	// Get all ways that share any of these nodes
	ways, err := s.WayRepository.GetWaysByNodeIDs(way.NodeIDs)
	if err != nil {
		return nil, err
	}

	// Filter out the original way
	adjacent := []models.Way{}
	for _, w := range ways {
		if w.ID != wayID {
			adjacent = append(adjacent, w)
		}
	}

	return adjacent, nil
}
