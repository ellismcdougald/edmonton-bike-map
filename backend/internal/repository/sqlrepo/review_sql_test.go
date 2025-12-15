package sqlrepo

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSQLReviewRepository_CreateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	rev := &models.Review{WayIDs: []int64{1}, UserID: 2, Rating: 5, Comment: "nice"}

	// Expect duplicate check before transaction
	mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT COUNT(1)
				FROM reviews r
				JOIN review_ways rw ON rw.review_id = r.id
				WHERE r.user_id = $1 AND rw.way_id = $2
			`)).
		WithArgs(int64(2), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect insert of review returning id
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(2), 5, "nice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	// Expect link insert into review_ways
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect transaction commit
	mock.ExpectCommit()

	repo := NewSQLReviewRepository(db)
	err = repo.CreateReview(rev)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_GetReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	created := time.Now()
	rows := sqlmock.NewRows([]string{"user_id", "rating", "comment", "created_at", "username"}).
		AddRow(int64(2), 4, "ok", created, "bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        JOIN review_ways rw ON rw.review_id = r.id\n        LEFT JOIN users u ON r.user_id = u.id\n        WHERE rw.way_id = $1;")).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	repo := NewSQLReviewRepository(db)
	reviews, err := repo.GetReviews(1)
	require.NoError(t, err)
	require.Len(t, reviews, 1)
	require.Equal(t, int64(2), reviews[0].UserID)
	require.Equal(t, "bob", reviews[0].Username)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_GetAllReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	created := time.Now()
	rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
		AddRow(int64(1), int64(2), 4, "ok", created, "bob").
		AddRow(int64(2), int64(3), 5, "great", created, "alice")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            rw.way_id,\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        JOIN review_ways rw ON rw.review_id = r.id\n        LEFT JOIN users u ON r.user_id = u.id;")).
		WillReturnRows(rows)

	repo := NewSQLReviewRepository(db)
	res, err := repo.GetAllReviews()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Len(t, res[int64(1)], 1)
	require.Len(t, res[int64(2)], 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_InsertBatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	repo := NewSQLReviewRepository(db)

	// normal batch of two reviews - expect two insert+link sequences with transactions
	// First review transaction
	mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT COUNT(1)
				FROM reviews r
				JOIN review_ways rw ON rw.review_id = r.id
				WHERE r.user_id = $1 AND rw.way_id = $2
			`)).
		WithArgs(int64(2), int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(2), 5, "great").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(100)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(100), int64(10)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Second review transaction
	mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT COUNT(1)
				FROM reviews r
				JOIN review_ways rw ON rw.review_id = r.id
				WHERE r.user_id = $1 AND rw.way_id = $2
			`)).
		WithArgs(int64(3), int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(3), 4, "ok").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(101), int64(11)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	reviews := []models.Review{{WayIDs: []int64{10}, UserID: 2, Rating: 5, Comment: "great"}, {WayIDs: []int64{11}, UserID: 3, Rating: 4, Comment: "ok"}}
	reqErr := repo.InsertBatches(reviews, 2)
	require.NoError(t, reqErr)

	// empty slice does nothing
	reqErr = repo.InsertBatches([]models.Review{}, 2)
	require.NoError(t, reqErr)

	// insert error on review insert
	mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT COUNT(1)
				FROM reviews r
				JOIN review_ways rw ON rw.review_id = r.id
				WHERE r.user_id = $1 AND rw.way_id = $2
			`)).
		WithArgs(int64(4), int64(12)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(4), 3, "meh").
		WillReturnError(fmt.Errorf("insert failure"))
	mock.ExpectRollback()

	reqErr = repo.InsertBatches([]models.Review{{WayIDs: []int64{12}, UserID: 4, Rating: 3, Comment: "meh"}}, 1)
	require.Error(t, reqErr)

	require.NoError(t, mock.ExpectationsWereMet())
}
