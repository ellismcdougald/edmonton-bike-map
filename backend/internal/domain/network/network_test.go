package network

import (
	"errors"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

// mock implementations of repository interfaces used by services
type mockNodeRepo struct {
	nodes map[int64]models.Node
	err   error
}

func (m *mockNodeRepo) Insert(node models.Node) error                          { return nil }
func (m *mockNodeRepo) InsertBatches(nodes []models.Node, batchSize int) error { return nil }
func (m *mockNodeRepo) GetNode(id int64) (*models.Node, error) {
	n, ok := m.nodes[id]
	if !ok {
		return nil, nil
	}
	return &n, nil
}
func (m *mockNodeRepo) GetAllNodes() (map[int64]models.Node, error) { return m.nodes, m.err }

type mockWayRepo struct {
	ways []models.Way
	err  error
}

func (m *mockWayRepo) Insert(way models.Way) error                          { return nil }
func (m *mockWayRepo) InsertBatches(ways []models.Way, batchSize int) error { return nil }
func (m *mockWayRepo) GetWay(id int64) (*models.Way, error) {
	for _, w := range m.ways {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, nil
}
func (m *mockWayRepo) GetAllWays() ([]models.Way, error) { return m.ways, m.err }
func (m *mockWayRepo) GetNearestWay(latitude, longitude float64) (*models.Way, error) {
	return nil, nil
}
func (m *mockWayRepo) GetWaysByNodeIDs(nodeIDs []int64) ([]models.Way, error) {
	return nil, nil
}

func TestBuildNetwork_SuccessAndEmpty(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
		2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
	}
	ways := []models.Way{{ID: 11, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}}

	nodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: nodes}}
	waySvc := service.WayService{WayRepository: &mockWayRepo{ways: ways}}

	net, err := BuildNetwork(nodeSvc, waySvc)
	require.NoError(t, err)
	require.Equal(t, nodes, net.Nodes)
	require.NotNil(t, net.Edges)

	// empty datasets -> no panic and empty network
	emptyNodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	emptyWaySvc := service.WayService{WayRepository: &mockWayRepo{ways: []models.Way{}}}

	net2, err := BuildNetwork(emptyNodeSvc, emptyWaySvc)
	require.NoError(t, err)
	require.NotNil(t, net2)
	require.Empty(t, net2.Nodes)
}

func TestBuildNetwork_PropagatesErrors(t *testing.T) {
	want := errors.New("node failure")
	nodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: nil, err: want}}
	waySvc := service.WayService{WayRepository: &mockWayRepo{ways: nil}}

	net, err := BuildNetwork(nodeSvc, waySvc)
	require.Nil(t, net)
	require.ErrorIs(t, err, want)

	want2 := errors.New("way failure")
	nodeSvc2 := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	waySvc2 := service.WayService{WayRepository: &mockWayRepo{ways: nil, err: want2}}
	net2, err2 := BuildNetwork(nodeSvc2, waySvc2)
	require.Nil(t, net2)
	require.ErrorIs(t, err2, want2)

	nodeSvc3 := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	waySvc3 := service.WayService{WayRepository: &mockWayRepo{ways: []models.Way{}}}
	net3, err3 := BuildNetwork(nodeSvc3, waySvc3)
	require.NoError(t, err3)
	require.NotNil(t, net3)
}
