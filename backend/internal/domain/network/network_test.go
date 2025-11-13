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

type mockReviewRepo struct {
	revs map[int64][]models.Review
	err  error
}

func (m *mockReviewRepo) CreateReview(review *models.Review) error                   { return nil }
func (m *mockReviewRepo) GetReviews(wayID int64) ([]models.Review, error)            { return m.revs[wayID], nil }
func (m *mockReviewRepo) GetAllReviews() (map[int64][]models.Review, error)          { return m.revs, m.err }
func (m *mockReviewRepo) InsertBatches(reviews []models.Review, batchSize int) error { return nil }

func TestBuildNetwork_SuccessAndEmpty(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
		2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
	}
	ways := []models.Way{{ID: 11, Tags: map[string]string{}, NodeIDs: []int64{1, 2}}}
	reviews := map[int64][]models.Review{}

	nodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: nodes}}
	waySvc := service.WayService{WayRepository: &mockWayRepo{ways: ways}}
	reviewSvc := service.ReviewService{ReviewRepository: &mockReviewRepo{revs: reviews}}

	net, err := BuildNetwork(nodeSvc, waySvc, reviewSvc)
	require.NoError(t, err)
	require.Equal(t, nodes, net.Nodes)
	require.NotNil(t, net.Edges)

	// empty datasets -> no panic and empty network
	emptyNodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	emptyWaySvc := service.WayService{WayRepository: &mockWayRepo{ways: []models.Way{}}}
	emptyReviewSvc := service.ReviewService{ReviewRepository: &mockReviewRepo{revs: map[int64][]models.Review{}}}

	net2, err := BuildNetwork(emptyNodeSvc, emptyWaySvc, emptyReviewSvc)
	require.NoError(t, err)
	require.NotNil(t, net2)
	require.Empty(t, net2.Nodes)
}

func TestBuildNetwork_PropagatesErrors(t *testing.T) {
	want := errors.New("node failure")
	nodeSvc := service.NodeService{NodeRepository: &mockNodeRepo{nodes: nil, err: want}}
	waySvc := service.WayService{WayRepository: &mockWayRepo{ways: nil}}
	reviewSvc := service.ReviewService{ReviewRepository: &mockReviewRepo{revs: nil}}

	net, err := BuildNetwork(nodeSvc, waySvc, reviewSvc)
	require.Nil(t, net)
	require.ErrorIs(t, err, want)

	want2 := errors.New("way failure")
	nodeSvc2 := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	waySvc2 := service.WayService{WayRepository: &mockWayRepo{ways: nil, err: want2}}
	reviewSvc2 := service.ReviewService{ReviewRepository: &mockReviewRepo{revs: nil}}
	net2, err2 := BuildNetwork(nodeSvc2, waySvc2, reviewSvc2)
	require.Nil(t, net2)
	require.ErrorIs(t, err2, want2)

	want3 := errors.New("review failure")
	nodeSvc3 := service.NodeService{NodeRepository: &mockNodeRepo{nodes: map[int64]models.Node{}}}
	waySvc3 := service.WayService{WayRepository: &mockWayRepo{ways: []models.Way{}}}
	reviewSvc3 := service.ReviewService{ReviewRepository: &mockReviewRepo{revs: nil, err: want3}}
	net3, err3 := BuildNetwork(nodeSvc3, waySvc3, reviewSvc3)
	require.Nil(t, net3)
	require.ErrorIs(t, err3, want3)
}
