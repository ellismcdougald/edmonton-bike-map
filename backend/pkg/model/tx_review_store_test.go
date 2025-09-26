package model

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestTxReviewStore_GetAllReviews(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(mock sqlmock.Sqlmock)
		wantReviews []Review
		wantErr     bool
	}{
		{
			name: "multiple reviews with usernames",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
					AddRow(int64(100), int64(1), 5, "Great bike path!", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), "user1").
					AddRow(int64(101), int64(2), 4, "Good route", time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC), "user2").
					AddRow(int64(102), int64(1), 3, "Okay path", time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC), "user1")
				mock.ExpectQuery(`SELECT
			r\.way_id,
			r\.user_id,
			r\.rating,
			r\.comment,
			r\.created_at,
			u\.username
		FROM reviews r
		LEFT JOIN users u ON r\.user_id = u\.id;`).
					WillReturnRows(rows)
			},
			wantReviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great bike path!", CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), Username: "user1"},
				{WayID: 101, UserID: 2, Rating: 4, Comment: "Good route", CreatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC), Username: "user2"},
				{WayID: 102, UserID: 1, Rating: 3, Comment: "Okay path", CreatedAt: time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC), Username: "user1"},
			},
			wantErr: false,
		},
		{
			name: "no reviews",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"})
				mock.ExpectQuery(`SELECT
			r\.way_id,
			r\.user_id,
			r\.rating,
			r\.comment,
			r\.created_at,
			u\.username
		FROM reviews r
		LEFT JOIN users u ON r\.user_id = u\.id;`).
					WillReturnRows(rows)
			},
			wantReviews: []Review{},
			wantErr:     false,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT
			r\.way_id,
			r\.user_id,
			r\.rating,
			r\.comment,
			r\.created_at,
			u\.username
		FROM reviews r
		LEFT JOIN users u ON r\.user_id = u\.id;`).
					WillReturnError(errors.New("query failed"))
			},
			wantReviews: []Review{},
			wantErr:     true,
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
					AddRow(int64(100), "invalid_user_id", 5, "Great!", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), "user1")
				mock.ExpectQuery(`SELECT
			r\.way_id,
			r\.user_id,
			r\.rating,
			r\.comment,
			r\.created_at,
			u\.username
		FROM reviews r
		LEFT JOIN users u ON r\.user_id = u\.id;`).
					WillReturnRows(rows)
			},
			wantReviews: []Review{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxReviewStore{Tx: tx}
			reviews, err := store.GetAllReviews()

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, reviews)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantReviews, reviews)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxReviewStore_GetReviews(t *testing.T) {
	tests := []struct {
		name        string
		wayID       int64
		mockSetup   func(mock sqlmock.Sqlmock)
		wantReviews []Review
		wantErr     bool
	}{
		{
			name:  "reviews found for way",
			wayID: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
					AddRow(int64(100), int64(1), 5, "Great bike path!", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), "user1").
					AddRow(int64(100), int64(2), 4, "Good route", time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC), "user2")
				mock.ExpectQuery(`SELECT\s+r\.way_id,\s+r\.user_id,.*?r\.rating,\s+r\.comment,\s+r\.created_at,\s+u\.username\s+FROM reviews r\s+LEFT JOIN users u ON r\.user_id = u\.id\s+WHERE r\.way_id = \$1`).
					WithArgs(int64(100)).
					WillReturnRows(rows)
			},
			wantReviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great bike path!", CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), Username: "user1"},
				{WayID: 100, UserID: 2, Rating: 4, Comment: "Good route", CreatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC), Username: "user2"},
			},
			wantErr: false,
		},
		{
			name:  "no reviews for way",
			wayID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"})
				mock.ExpectQuery(`SELECT\s+r\.way_id,\s+r\.user_id,.*?r\.rating,\s+r\.comment,\s+r\.created_at,\s+u\.username\s+FROM reviews r\s+LEFT JOIN users u ON r\.user_id = u\.id\s+WHERE r\.way_id = \$1`).
					WithArgs(int64(999)).
					WillReturnRows(rows)
			},
			wantReviews: []Review{},
			wantErr:     false,
		},
		{
			name:  "query error",
			wayID: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT\s+r\.way_id,\s+r\.user_id,.*?r\.rating,\s+r\.comment,\s+r\.created_at,\s+u\.username\s+FROM reviews r\s+LEFT JOIN users u ON r\.user_id = u\.id\s+WHERE r\.way_id = \$1`).
					WithArgs(int64(100)).
					WillReturnError(errors.New("query failed"))
			},
			wantReviews: nil,
			wantErr:     true,
		},
		{
			name:  "scan error",
			wayID: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
					AddRow(int64(100), "invalid_user_id", 5, "Great!", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), "user1")
				mock.ExpectQuery(`SELECT\s+r\.way_id,\s+r\.user_id,.*?r\.rating,\s+r\.comment,\s+r\.created_at,\s+u\.username\s+FROM reviews r\s+LEFT JOIN users u ON r\.user_id = u\.id\s+WHERE r\.way_id = \$1`).
					WithArgs(int64(100)).
					WillReturnRows(rows)
			},
			wantReviews: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxReviewStore{Tx: tx}
			reviews, err := store.GetReviews(tt.wayID)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, reviews)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantReviews, reviews)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxReviewStore_InsertReview(t *testing.T) {
	tests := []struct {
		name      string
		review    *Review
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful insert",
			review: &Review{
				WayID:   100,
				UserID:  1,
				Rating:  5,
				Comment: "Great bike path!",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(
			way_id,
			user_id,
			rating,
			comment
		\)
		VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(100), int64(1), 5, "Great bike path!").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "insert error",
			review: &Review{
				WayID:   101,
				UserID:  2,
				Rating:  4,
				Comment: "Good route",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(
			way_id,
			user_id,
			rating,
			comment
		\)
		VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(101), int64(2), 4, "Good route").
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxReviewStore{Tx: tx}
			err = store.InsertReview(tt.review)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxReviewStore_InsertBatch(t *testing.T) {
	tests := []struct {
		name      string
		reviews   []Review
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful batch insert",
			reviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"},
				{WayID: 101, UserID: 2, Rating: 4, Comment: "Good"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\), \(\$5, \$6, \$7, \$8\)`).
					WithArgs(int64(100), int64(1), 5, "Great!", int64(101), int64(2), 4, "Good").
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: false,
		},
		{
			name: "single review batch",
			reviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(100), int64(1), 5, "Great!").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:    "empty slice does nothing",
			reviews: []Review{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: false,
		},
		{
			name: "insert error",
			reviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(100), int64(1), 5, "Great!").
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxReviewStore{Tx: tx}
			err = store.InsertBatch(tt.reviews)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxReviewStore_InsertBatchChunks(t *testing.T) {
	tests := []struct {
		name      string
		reviews   []Review
		batchSize int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "multiple chunks",
			reviews: []Review{
				{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"},
				{WayID: 101, UserID: 2, Rating: 4, Comment: "Good"},
				{WayID: 102, UserID: 3, Rating: 3, Comment: "Okay"},
			},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// First batch (2 reviews)
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\), \(\$5, \$6, \$7, \$8\)`).
					WithArgs(int64(100), int64(1), 5, "Great!", int64(101), int64(2), 4, "Good").
					WillReturnResult(sqlmock.NewResult(2, 2))
				// Second batch (1 review)
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(102), int64(3), 3, "Okay").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:      "default batch size when <= 0",
			reviews:   []Review{{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"}},
			batchSize: 0, // Should default to 1000
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(100), int64(1), 5, "Great!").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:      "empty slice does nothing",
			reviews:   []Review{},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: false,
		},
		{
			name:      "insert batch error propagates",
			reviews:   []Review{{WayID: 100, UserID: 1, Rating: 5, Comment: "Great!"}},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO reviews \(way_id, user_id, rating, comment\) VALUES \(\$1, \$2, \$3, \$4\)`).
					WithArgs(int64(100), int64(1), 5, "Great!").
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxReviewStore{Tx: tx}
			err = store.InsertBatchChunks(tt.reviews, tt.batchSize)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}