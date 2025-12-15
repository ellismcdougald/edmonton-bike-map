package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

// mockReviewRepo implements minimal ReviewRepository for tests.
type mockReviewRepo struct {
	reviews map[int64][]models.Review
	err     error
	created []*models.Review
}

func (m *mockReviewRepo) CreateReview(review *models.Review) error {
	if m.err != nil {
		return m.err
	}
	m.created = append(m.created, review)
	// support multi-way by adding review to each way ID
	wayIDs := review.WayIDs
	// Do not rely on deprecated WayID in tests; ensure WayIDs carries target ways
	for _, wid := range wayIDs {
		m.reviews[wid] = append(m.reviews[wid], *review)
	}
	return nil
}
func (m *mockReviewRepo) GetReviews(wayID int64) ([]models.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reviews[wayID], nil
}
func (m *mockReviewRepo) GetAllReviews() (map[int64][]models.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reviews, nil
}
func (m *mockReviewRepo) InsertBatches(reviews []models.Review, batchSize int) error {
	if m.err != nil {
		return m.err
	}
	for _, r := range reviews {
		for _, wid := range r.WayIDs {
			m.reviews[wid] = append(m.reviews[wid], r)
		}
	}
	return nil
}
func (m *mockReviewRepo) DeleteUserReviewForWay(userID int64, wayID int64) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestReviewHandler_HandleGetReviews_Success(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{
		1: {{WayIDs: []int64{1}, UserID: 2, Rating: 4, Comment: "ok", Username: "bob"}},
	}}
	svc := service.NewReviewService(repo)
	h := NewReviewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reviews?wayID=1", nil)
	rr := httptest.NewRecorder()

	handler := h.HandleGetReviews()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var got []models.Review
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Len(t, got, 1)
	// Ensure review is associated to requested way via repository grouping
	// WayIDs not returned in payload; we only validate the presence and content
}

func TestReviewHandler_HandleGetReviews_BadRequest(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := service.NewReviewService(repo)
	h := NewReviewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reviews?wayID=notanint", nil)
	rr := httptest.NewRecorder()

	handler := h.HandleGetReviews()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestReviewHandler_HandlePostReview_Success(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := service.NewReviewService(repo)
	h := NewReviewHandler(svc)

	// new payload uses wayId and omits userId (taken from context)
	payload := map[string]interface{}{"wayId": 1, "rating": 5, "comment": "nice"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
	// Set user ID in context using middleware helper
	ctx := middleware.UserIDToContext(req.Context(), int64(2))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := h.HandlePostReview()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	// Ensure repo received the review
	require.Len(t, repo.created, 1)
	require.Equal(t, []int64{1}, repo.created[0].WayIDs)
}

func TestReviewHandler_HandlePostReview_Unauthorized(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := service.NewReviewService(repo)
	h := NewReviewHandler(svc)

	payload := map[string]interface{}{"wayId": 1, "rating": 5, "comment": "nice"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
	// Don't set user ID in context - should fail with 401
	rr := httptest.NewRecorder()

	handler := h.HandlePostReview()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestReviewHandler_HandlePostReview_BadRating(t *testing.T) {
	repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
	svc := service.NewReviewService(repo)
	h := NewReviewHandler(svc)

	payload := map[string]interface{}{"wayId": 1, "rating": 0, "comment": "bad"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
	// Set user ID in context using middleware helper
	ctx := middleware.UserIDToContext(req.Context(), int64(2))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := h.HandlePostReview()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
