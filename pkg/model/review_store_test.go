package model_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &model.DBReviewStore{DB: db}

	tests := []struct {
		name      string
		review    *model.Review
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful insert",
			review: &model.Review{
				WayID:   1,
				UserID:  2,
				Rating:  4,
				Comment: "Nice trail",
			},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO reviews`).
					WithArgs(int64(1), int64(2), 4, "Nice trail").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "insert failure",
			review: &model.Review{
				WayID:   1,
				UserID:  2,
				Rating:  3,
				Comment: "Okay trail",
			},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO reviews`).
					WithArgs(int64(1), int64(2), 3, "Okay trail").
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := store.CreateReview(tt.review)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &model.DBReviewStore{DB: db}

	columns := []string{"way_id", "rating", "comment", "created_at", "username"}

	createdAt1 := time.Now()
	createdAt2 := createdAt1.Add(time.Hour)

	tests := []struct {
		name      string
		wayID     int64
		mockSetup func()
		want      []model.Review
		wantErr   bool
	}{
		{
			name:  "multiple reviews found",
			wayID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows(columns).
					AddRow(int64(1), 5, "Great!", createdAt1, "user1").
					AddRow(int64(1), 3, "Not bad", createdAt2, "user2")
				mock.ExpectQuery(`SELECT`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			want: []model.Review{
				{WayID: 1, Rating: 5, Comment: "Great!", CreatedAt: createdAt1, Username: "user1"},
				{WayID: 1, Rating: 3, Comment: "Not bad", CreatedAt: createdAt2, Username: "user2"},
			},
			wantErr: false,
		},
		{
			name:  "no reviews found",
			wayID: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(`SELECT`).
					WithArgs(int64(2)).
					WillReturnRows(rows)
			},
			want:    []model.Review{},
			wantErr: false,
		},
		{
			name:  "query failure",
			wayID: 3,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).
					WithArgs(int64(3)).
					WillReturnError(errors.New("query failed"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "scan failure",
			wayID: 4,
			mockSetup: func() {
				rows := sqlmock.NewRows(columns).
					AddRow("bad_way_id", 5, "Invalid data", createdAt1, "user1")
				mock.ExpectQuery(`SELECT`).
					WithArgs(int64(4)).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := store.GetReviews(tt.wayID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
