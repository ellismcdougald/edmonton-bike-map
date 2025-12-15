package service

import (
	"errors"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

type mockReviewRepo struct {
	created []models.Review
	deleted []struct {
		userID int64
		wayID  int64
	}
	reviews map[int64][]models.Review
	err     error
}

func (m *mockReviewRepo) CreateReview(review *models.Review) error {
	m.created = append(m.created, *review)
	return m.err
}

func (m *mockReviewRepo) GetReviews(wayID int64) ([]models.Review, error) {
	return m.reviews[wayID], m.err
}

func (m *mockReviewRepo) GetAllReviews() (map[int64][]models.Review, error) {
	return m.reviews, m.err
}

func (m *mockReviewRepo) InsertBatches(reviews []models.Review, batchSize int) error {
	return m.err
}

func (m *mockReviewRepo) DeleteUserReviewForWay(userID int64, wayID int64) error {
	m.deleted = append(m.deleted, struct {
		userID int64
		wayID  int64
	}{userID, wayID})
	return m.err
}

func TestReviewService_AddReview_DelegatesAndPassesData(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := NewReviewService(repo)

	review := models.Review{WayIDs: []int64{10}, UserID: 1, Rating: 4, Comment: "great"}

	err := svc.AddReview(review)
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	require.Equal(t, review, repo.created[0])
}

func TestReviewService_AddReview_PropagatesError(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}, err: errors.New("create fail")}
	svc := NewReviewService(repo)

	err := svc.AddReview(models.Review{WayIDs: []int64{1}, UserID: 2, Rating: 5})
	require.ErrorIs(t, err, repo.err)
}

func TestReviewService_GetReviewsForWay_ReturnsData(t *testing.T) {
	expected := []models.Review{{UserID: 3, Rating: 5, Comment: "smooth"}}
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{42: expected}}
	svc := NewReviewService(repo)

	res, err := svc.GetReviewsForWay(42)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestReviewService_GetReviewsForWay_PropagatesError(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}, err: errors.New("get fail")}
	svc := NewReviewService(repo)

	_, err := svc.GetReviewsForWay(11)
	require.ErrorIs(t, err, repo.err)
}

func TestReviewService_GetAllReviews_ReturnsData(t *testing.T) {
	data := map[int64][]models.Review{
		5: {{UserID: 1, Rating: 4}},
		6: {{UserID: 2, Rating: 5}},
	}
	repo := &mockReviewRepo{reviews: data}
	svc := NewReviewService(repo)

	res, err := svc.GetAllReviews()
	require.NoError(t, err)
	require.Equal(t, data, res)
}

func TestReviewService_GetAllReviews_PropagatesError(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}, err: errors.New("all fail")}
	svc := NewReviewService(repo)

	_, err := svc.GetAllReviews()
	require.ErrorIs(t, err, repo.err)
}

func TestReviewService_DeleteUserReviewForWay_Delegates(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := NewReviewService(repo)

	err := svc.DeleteUserReviewForWay(9, 44)
	require.NoError(t, err)
	require.Len(t, repo.deleted, 1)
	require.Equal(t, int64(9), repo.deleted[0].userID)
	require.Equal(t, int64(44), repo.deleted[0].wayID)
}

func TestReviewService_DeleteUserReviewForWay_PropagatesError(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}, err: errors.New("delete fail")}
	svc := NewReviewService(repo)

	err := svc.DeleteUserReviewForWay(3, 8)
	require.ErrorIs(t, err, repo.err)
	require.Len(t, repo.deleted, 1)
	require.Equal(t, int64(3), repo.deleted[0].userID)
	require.Equal(t, int64(8), repo.deleted[0].wayID)
}
