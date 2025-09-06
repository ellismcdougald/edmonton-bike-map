package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
)

type mockReviewService struct {
	GetReviewsFunc    func(wayID int64) ([]model.Review, error)
	CreateReviewFunc  func(review *model.Review) error
	GetAllReviewsFunc func() (map[int][]model.Review, error)
}

func (m *mockReviewService) GetReviews(wayID int64) ([]model.Review, error) {
	return m.GetReviewsFunc(wayID)
}

func (m *mockReviewService) CreateReview(review *model.Review) error {
	return m.CreateReviewFunc(review)
}

func (m *mockReviewService) GetAllReviews() (map[int][]model.Review, error) {
	return m.GetAllReviewsFunc()
}

func TestHandleGetReviews(t *testing.T) {
	tests := []struct {
		name             string
		wayID            string
		mockServiceFunc  func(wayID int64) ([]model.Review, error)
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:  "valid wayID",
			wayID: "123",
			mockServiceFunc: func(wayID int64) ([]model.Review, error) {
				return []model.Review{
					{
						WayID:     123,
						UserID:    1,
						Rating:    8,
						Comment:   "Nice route",
						CreatedAt: time.Now(),
						Username:  "bob",
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "invalid wayID",
			wayID: "abc",
			mockServiceFunc: func(wayID int64) ([]model.Review, error) {
				return nil, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "db error",
			wayID: "999",
			mockServiceFunc: func(wayID int64) ([]model.Review, error) {
				return nil, errors.New("db failure")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := handlers.RealHandlers{
				ReviewService: &mockReviewService{
					GetReviewsFunc: tc.mockServiceFunc,
				},
			}

			req := httptest.NewRequest("GET", "/api/reviews?wayID="+tc.wayID, nil)
			w := httptest.NewRecorder()
			handler.HandleGetReviews().ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandlePostReview(t *testing.T) {
	tests := []struct {
		name            string
		body            model.Review
		mockServiceFunc func(review *model.Review) error
		expectedStatus  int
	}{
		{
			name: "valid review",
			body: model.Review{WayID: 123, UserID: 1, Rating: 8, Comment: "Great!", Username: "bob"},
			mockServiceFunc: func(review *model.Review) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "bad json",
			body: model.Review{},
			mockServiceFunc: func(review *model.Review) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid rating",
			body: model.Review{WayID: 123, UserID: 1, Rating: 20, Comment: "oops"},
			mockServiceFunc: func(review *model.Review) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: model.Review{WayID: 123, UserID: 1, Rating: 6, Comment: "ok"},
			mockServiceFunc: func(review *model.Review) error {
				return errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := handlers.RealHandlers{
				ReviewService: &mockReviewService{
					CreateReviewFunc: tc.mockServiceFunc,
				},
			}

			var bodyBytes []byte
			if tc.name == "bad json" {
				bodyBytes = []byte(`{"invalid_json"`)
			} else {
				bodyBytes, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest("POST", "/api/reviews", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandlePostReview().ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
