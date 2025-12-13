package service

import (
	"errors"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
	"github.com/stretchr/testify/require"
)

// mockWayRepoForService implements WayRepository for service-layer testing
type mockWayRepoForService struct {
	ways         []models.Way
	waysByNodeID []models.Way
	getWayErr    error
	getByNodeErr error
}

func (m *mockWayRepoForService) Insert(way models.Way) error { return nil }
func (m *mockWayRepoForService) InsertBatches(ways []models.Way, batchSize int) error {
	return nil
}
func (m *mockWayRepoForService) GetWay(id int64) (*models.Way, error) {
	if m.getWayErr != nil {
		return nil, m.getWayErr
	}
	for _, w := range m.ways {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, repository.ErrWayNotFound
}
func (m *mockWayRepoForService) GetAllWays() ([]models.Way, error) {
	return m.ways, nil
}
func (m *mockWayRepoForService) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	return nil, nil
}
func (m *mockWayRepoForService) GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error) {
	if m.getByNodeErr != nil {
		return nil, m.getByNodeErr
	}
	return m.waysByNodeID, nil
}

func TestWayService_GetAdjacentWays(t *testing.T) {
	t.Run("returns adjacent ways excluding target way", func(t *testing.T) {
		targetWay := models.Way{ID: 10, Tags: map[string]string{"name": "Main"}, NodeIDs: []int64{1, 2}}
		adjacentWay1 := models.Way{ID: 20, Tags: map[string]string{"name": "Side"}, NodeIDs: []int64{2, 3}}
		adjacentWay2 := models.Way{ID: 30, Tags: map[string]string{"name": "Cross"}, NodeIDs: []int64{1, 4}}

		mockRepo := &mockWayRepoForService{
			ways:         []models.Way{targetWay, adjacentWay1, adjacentWay2},
			waysByNodeID: []models.Way{targetWay, adjacentWay1, adjacentWay2}, // Simulates DB returning all ways with shared nodes
		}

		svc := NewWayService(mockRepo)
		adjacent, err := svc.GetAdjacentWays(10)

		require.NoError(t, err)
		require.Len(t, adjacent, 2)
		require.Equal(t, int64(20), adjacent[0].ID)
		require.Equal(t, int64(30), adjacent[1].ID)
	})

	t.Run("returns empty slice when no adjacent ways", func(t *testing.T) {
		targetWay := models.Way{ID: 10, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}

		mockRepo := &mockWayRepoForService{
			ways:         []models.Way{targetWay},
			waysByNodeID: []models.Way{targetWay}, // Only the target way shares its own nodes
		}

		svc := NewWayService(mockRepo)
		adjacent, err := svc.GetAdjacentWays(10)

		require.NoError(t, err)
		require.Empty(t, adjacent)
	})

	t.Run("returns empty slice when way has no nodes", func(t *testing.T) {
		targetWay := models.Way{ID: 10, Tags: map[string]string{}, NodeIDs: []int64{}}

		mockRepo := &mockWayRepoForService{
			ways:         []models.Way{targetWay},
			waysByNodeID: []models.Way{},
		}

		svc := NewWayService(mockRepo)
		adjacent, err := svc.GetAdjacentWays(10)

		require.NoError(t, err)
		require.Empty(t, adjacent)
	})

	t.Run("returns error when way not found", func(t *testing.T) {
		mockRepo := &mockWayRepoForService{
			ways:      []models.Way{},
			getWayErr: repository.ErrWayNotFound,
		}

		svc := NewWayService(mockRepo)
		adjacent, err := svc.GetAdjacentWays(999)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrWayNotFound)
		require.Nil(t, adjacent)
	})

	t.Run("propagates GetWaysByNodeIDs error", func(t *testing.T) {
		targetWay := models.Way{ID: 10, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}
		dbErr := repository.ErrWayNotFound

		mockRepo := &mockWayRepoForService{
			ways:         []models.Way{targetWay},
			getByNodeErr: dbErr,
		}

		svc := NewWayService(mockRepo)
		adjacent, err := svc.GetAdjacentWays(10)

		require.Error(t, err)
		require.True(t, errors.Is(err, dbErr))
		require.Nil(t, adjacent)
	})
}
