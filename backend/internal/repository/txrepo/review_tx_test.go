package txrepo

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestTxReviewRepository_CreateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	rev := &models.Review{WayIDs: []int64{1}, UserID: 2, Rating: 5, Comment: "nice"}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(2), 5, "nice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(200)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(200), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewTxReviewRepository(tx)
	err = repo.CreateReview(rev)
	require.NoError(t, err)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxReviewRepository_DeleteUserReviewForWay(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := &TxReviewRepository{Tx: tx}

	find := `
		SELECT r.id
		FROM reviews r
		JOIN review_ways rw ON rw.review_id = r.id
		WHERE r.user_id = $1 AND rw.way_id = $2
		LIMIT 1
	`

	// happy path
	mock.ExpectQuery(regexp.QuoteMeta(find)).
		WithArgs(int64(5), int64(12)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(77)))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM review_ways WHERE review_id = $1")).
		WithArgs(int64(77)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM reviews WHERE id = $1")).
		WithArgs(int64(77)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	require.NoError(t, repo.DeleteUserReviewForWay(5, 12))

	// no rows path
	mock.ExpectQuery(regexp.QuoteMeta(find)).
		WithArgs(int64(6), int64(99)).
		WillReturnError(sql.ErrNoRows)

	require.NoError(t, repo.DeleteUserReviewForWay(6, 99))

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxReviewRepository_InsertBatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := NewTxReviewRepository(tx)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(5), 5, "super").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(201)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(201), int64(20)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(6), 4, "nice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(202)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)")).
		WithArgs(int64(202), int64(21)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	reviews := []models.Review{{WayIDs: []int64{20}, UserID: 5, Rating: 5, Comment: "super"}, {WayIDs: []int64{21}, UserID: 6, Rating: 4, Comment: "nice"}}
	reqErr := repo.InsertBatches(reviews, 2)
	require.NoError(t, reqErr)

	// empty slice
	reqErr = repo.InsertBatches([]models.Review{}, 2)
	require.NoError(t, reqErr)

	// exec error
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO reviews (\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3)\n        RETURNING id")).
		WithArgs(int64(7), 2, "bad").
		WillReturnError(fmt.Errorf("insert failure"))

	reqErr = repo.InsertBatches([]models.Review{{WayIDs: []int64{22}, UserID: 7, Rating: 2, Comment: "bad"}}, 1)
	require.Error(t, reqErr)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxReviewRepository_GetReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	created := time.Now()
	rows := sqlmock.NewRows([]string{"user_id", "rating", "comment", "created_at", "username"}).
		AddRow(int64(2), 4, "ok", created, "bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        JOIN review_ways rw ON rw.review_id = r.id\n        LEFT JOIN users u ON r.user_id = u.id\n        WHERE rw.way_id = $1;")).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	repo := NewTxReviewRepository(tx)
	reviews, err := repo.GetReviews(1)
	require.NoError(t, err)
	require.Len(t, reviews, 1)
	require.Equal(t, int64(2), reviews[0].UserID)
	require.Equal(t, "bob", reviews[0].Username)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxReviewRepository_GetAllReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	created := time.Now()
	rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
		AddRow(int64(1), int64(2), 4, "ok", created, "bob").
		AddRow(int64(2), int64(3), 5, "great", created, "alice")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            rw.way_id,\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        JOIN review_ways rw ON rw.review_id = r.id\n        LEFT JOIN users u ON r.user_id = u.id;")).
		WillReturnRows(rows)

	repo := NewTxReviewRepository(tx)
	res, err := repo.GetAllReviews()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Len(t, res[int64(1)], 1)
	require.Len(t, res[int64(2)], 1)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}
