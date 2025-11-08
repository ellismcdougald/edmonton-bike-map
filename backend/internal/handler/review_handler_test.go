package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/stretchr/testify/require"
)

// mockReviewRepo implements minimal ReviewRepository for tests.
type mockReviewRepo struct{
    reviews map[int64][]models.Review
    err     error
    created []*models.Review
}

func (m *mockReviewRepo) CreateReview(review *models.Review) error {
    if m.err != nil { return m.err }
    m.created = append(m.created, review)
    m.reviews[review.WayID] = append(m.reviews[review.WayID], *review)
    return nil
}
func (m *mockReviewRepo) GetReviews(wayID int64) ([]models.Review, error) {
    if m.err != nil { return nil, m.err }
    return m.reviews[wayID], nil
}
func (m *mockReviewRepo) GetAllReviews() (map[int64][]models.Review, error) {
    if m.err != nil { return nil, m.err }
    return m.reviews, nil
}

func TestReviewHandler_HandleGetReviews_Success(t *testing.T) {
    repo := &mockReviewRepo{reviews: map[int64][]models.Review{
        1: {{WayID: 1, UserID: 2, Rating: 4, Comment: "ok", Username: "bob"}},
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
    require.Equal(t, int64(1), got[0].WayID)
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

    body := models.Review{WayID: 1, UserID: 2, Rating: 5, Comment: "nice"}
    b, _ := json.Marshal(body)

    req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
    rr := httptest.NewRecorder()

    handler := h.HandlePostReview()
    handler.ServeHTTP(rr, req)

    require.Equal(t, http.StatusCreated, rr.Code)
    // Ensure repo received the review
    require.Len(t, repo.created, 1)
    require.Equal(t, int64(1), repo.created[0].WayID)
}

func TestReviewHandler_HandlePostReview_BadRating(t *testing.T) {
    repo := &mockReviewRepo{reviews: map[int64][]models.Review{}}
    svc := service.NewReviewService(repo)
    h := NewReviewHandler(svc)

    body := models.Review{WayID: 1, UserID: 2, Rating: 0, Comment: "bad"}
    b, _ := json.Marshal(body)

    req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
    rr := httptest.NewRecorder()

    handler := h.HandlePostReview()
    handler.ServeHTTP(rr, req)

    require.Equal(t, http.StatusBadRequest, rr.Code)
}
