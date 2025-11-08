package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

// mockNodeRepo implements the minimal NodeRepository for tests.
type mockNodeRepo struct{
    nodes map[int64]models.Node
    err   error
}

func (m *mockNodeRepo) Insert(node models.Node) error { return nil }
func (m *mockNodeRepo) InsertBatches(nodes []models.Node, batchSize int) error { return nil }
func (m *mockNodeRepo) GetNode(id int64) (*models.Node, error) {
    if m.err != nil { return nil, m.err }
    n, ok := m.nodes[id]
    if !ok { return nil, nil }
    return &n, nil
}
func (m *mockNodeRepo) GetAllNodes() (map[int64]models.Node, error) {
    if m.err != nil { return nil, m.err }
    return m.nodes, nil
}

// mockWayRepo implements the minimal WayRepository for tests.
type mockWayRepo struct{
    ways []models.Way
    err  error
}

func (m *mockWayRepo) Insert(way models.Way) error { return nil }
func (m *mockWayRepo) InsertBatches(ways []models.Way, batchSize int) error { return nil }
func (m *mockWayRepo) GetWay(id int64) (*models.Way, error) { return nil, nil }
func (m *mockWayRepo) GetAllWays() ([]models.Way, error) {
    if m.err != nil { return nil, m.err }
    return m.ways, nil
}

func TestWayHandler_HandleAllWays_Success(t *testing.T) {
    nodes := map[int64]models.Node{
        10: {ID: 10, Latitude: 1.0, Longitude: 2.0},
    }

    ways := []models.Way{{ID: 1, Tags: map[string]string{"a":"b"}, NodeIDs: []int64{10}}}

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
    ways := []models.Way{{ID: 1, Tags: map[string]string{"a":"b"}, NodeIDs: []int64{99}}}
    nodeService := service.NewNodeService(&mockNodeRepo{nodes: map[int64]models.Node{}})
    wayService := service.NewWayService(&mockWayRepo{ways: ways})

    h := NewWayHandler(nodeService, wayService)

    req := httptest.NewRequest(http.MethodGet, "/ways", nil)
    rr := httptest.NewRecorder()

    handler := h.HandleAllWays()
    handler.ServeHTTP(rr, req)

    require.Equal(t, http.StatusInternalServerError, rr.Code)
}
