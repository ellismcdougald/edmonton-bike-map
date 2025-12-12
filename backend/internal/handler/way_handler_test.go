package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

// mockNodeRepo implements the minimal NodeRepository for tests.
type mockNodeRepo struct {
	nodes map[int64]models.Node
	err   error
}

func (m *mockNodeRepo) Insert(node models.Node) error                          { return nil }
func (m *mockNodeRepo) InsertBatches(nodes []models.Node, batchSize int) error { return nil }
func (m *mockNodeRepo) GetNode(id int64) (*models.Node, error) {
	if m.err != nil {
		return nil, m.err
	}
	n, ok := m.nodes[id]
	if !ok {
		return nil, nil
	}
	return &n, nil
}
func (m *mockNodeRepo) GetAllNodes() (map[int64]models.Node, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.nodes, nil
}

// mockWayRepo implements the minimal WayRepository for tests.
type mockWayRepo struct {
	ways []models.Way
	err  error
}

func (m *mockWayRepo) Insert(way models.Way) error                          { return nil }
func (m *mockWayRepo) InsertBatches(ways []models.Way, batchSize int) error { return nil }
func (m *mockWayRepo) GetWay(id int64) (*models.Way, error)                 { return nil, nil }
func (m *mockWayRepo) GetAllWays() ([]models.Way, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.ways, nil
}
func (m *mockWayRepo) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	return nil, nil
}
func (m *mockWayRepo) GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error) {
	return nil, nil
}

func TestWayHandler_HandleAllWays_Success(t *testing.T) {
	nodes := map[int64]models.Node{
		10: {ID: 10, Latitude: 1.0, Longitude: 2.0},
	}

	ways := []models.Way{{ID: 1, Tags: map[string]string{"a": "b"}, NodeIDs: []int64{10}}}

	nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
	wayService := service.NewWayService(&mockWayRepo{ways: ways})

	h := NewWayHandler(nodeService, wayService)

	req := httptest.NewRequest(http.MethodGet, "/ways", nil)
	rr := httptest.NewRecorder()

	handler := h.HandleAllWays()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var fc models.FeatureCollection
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &fc))
	require.Len(t, fc.Features, 1)
	require.Equal(t, "1", fc.Features[0].Properties["id"])
}

func TestWayHandler_HandleAllWays_NodeMissing_MapsError(t *testing.T) {
	// Way references a node that doesn't exist in the nodes map -> mapping error -> 500
	ways := []models.Way{{ID: 1, Tags: map[string]string{"a": "b"}, NodeIDs: []int64{99}}}
	nodeService := service.NewNodeService(&mockNodeRepo{nodes: map[int64]models.Node{}})
	wayService := service.NewWayService(&mockWayRepo{ways: ways})

	h := NewWayHandler(nodeService, wayService)

	req := httptest.NewRequest(http.MethodGet, "/ways", nil)
	rr := httptest.NewRecorder()

	handler := h.HandleAllWays()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestWayHandler_HandleGetAdjacentWays(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
		2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
		3: {ID: 3, Latitude: 50.2, Longitude: -113.2},
	}

	t.Run("returns adjacent ways as GeoJSON FeatureCollection", func(t *testing.T) {
		// Way 10 has nodes [1, 2]
		// Way 20 shares node 2 and has nodes [2, 3]
		// Way 30 shares node 3 and has nodes [3]
		targetWay := models.Way{ID: 10, Tags: map[string]string{"name": "Main"}, NodeIDs: []int64{1, 2}}
		adjacentWays := []models.Way{
			{ID: 20, Tags: map[string]string{"name": "Side"}, NodeIDs: []int64{2, 3}},
			{ID: 30, Tags: map[string]string{"name": "Connector"}, NodeIDs: []int64{3}},
		}

		mockRepo := &mockWayRepoWithAdjacent{
			targetWay:    &targetWay,
			adjacentWays: adjacentWays,
		}

		nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
		wayService := service.NewWayService(mockRepo)
		h := NewWayHandler(nodeService, wayService)

		req := httptest.NewRequest(http.MethodGet, "/api/adjacent-ways?id=10", nil)
		rr := httptest.NewRecorder()

		handler := h.HandleGetAdjacentWays()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var fc models.FeatureCollection
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &fc))
		require.Equal(t, "FeatureCollection", fc.Type)
		require.Len(t, fc.Features, 2)
		require.Equal(t, "20", fc.Features[0].Properties["id"])
		require.Equal(t, "30", fc.Features[1].Properties["id"])
	})

	t.Run("returns 400 when id parameter missing", func(t *testing.T) {
		mockRepo := &mockWayRepoWithAdjacent{}
		nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
		wayService := service.NewWayService(mockRepo)
		h := NewWayHandler(nodeService, wayService)

		req := httptest.NewRequest(http.MethodGet, "/api/adjacent-ways", nil)
		rr := httptest.NewRecorder()

		handler := h.HandleGetAdjacentWays()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Contains(t, rr.Body.String(), "Missing required parameter")
	})

	t.Run("returns 400 when id parameter invalid", func(t *testing.T) {
		mockRepo := &mockWayRepoWithAdjacent{}
		nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
		wayService := service.NewWayService(mockRepo)
		h := NewWayHandler(nodeService, wayService)

		req := httptest.NewRequest(http.MethodGet, "/api/adjacent-ways?id=invalid", nil)
		rr := httptest.NewRecorder()

		handler := h.HandleGetAdjacentWays()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Contains(t, rr.Body.String(), "Invalid id")
	})

	t.Run("returns 404 when way not found", func(t *testing.T) {
		mockRepo := &mockWayRepoWithAdjacent{
			notFound: true,
		}
		nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
		wayService := service.NewWayService(mockRepo)
		h := NewWayHandler(nodeService, wayService)

		req := httptest.NewRequest(http.MethodGet, "/api/adjacent-ways?id=999", nil)
		rr := httptest.NewRecorder()

		handler := h.HandleGetAdjacentWays()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("returns empty FeatureCollection when no adjacent ways", func(t *testing.T) {
		targetWay := models.Way{ID: 10, Tags: map[string]string{}, NodeIDs: []int64{1}}
		mockRepo := &mockWayRepoWithAdjacent{
			targetWay:    &targetWay,
			adjacentWays: []models.Way{},
		}

		nodeService := service.NewNodeService(&mockNodeRepo{nodes: nodes})
		wayService := service.NewWayService(mockRepo)
		h := NewWayHandler(nodeService, wayService)

		req := httptest.NewRequest(http.MethodGet, "/api/adjacent-ways?id=10", nil)
		rr := httptest.NewRecorder()

		handler := h.HandleGetAdjacentWays()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var fc models.FeatureCollection
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &fc))
		require.Equal(t, "FeatureCollection", fc.Type)
		require.Empty(t, fc.Features)
	})
}

// mockWayRepoWithAdjacent provides controlled responses for adjacent way testing
type mockWayRepoWithAdjacent struct {
	targetWay    *models.Way
	adjacentWays []models.Way
	notFound     bool
	err          error
}

func (m *mockWayRepoWithAdjacent) Insert(way models.Way) error { return nil }
func (m *mockWayRepoWithAdjacent) InsertBatches(ways []models.Way, batchSize int) error {
	return nil
}
func (m *mockWayRepoWithAdjacent) GetWay(id int64) (*models.Way, error) {
	if m.notFound {
		return nil, repository.ErrWayNotFound
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.targetWay, nil
}
func (m *mockWayRepoWithAdjacent) GetAllWays() ([]models.Way, error) {
	return nil, nil
}
func (m *mockWayRepoWithAdjacent) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	return nil, nil
}
func (m *mockWayRepoWithAdjacent) GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Return target way + adjacent ways, simulating database behavior
	allWays := []models.Way{}
	if m.targetWay != nil {
		allWays = append(allWays, *m.targetWay)
	}
	allWays = append(allWays, m.adjacentWays...)
	return allWays, nil
}
