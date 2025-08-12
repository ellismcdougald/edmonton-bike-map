package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBUserStore_GetUser(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		mockSetup func(mock sqlmock.Sqlmock)
		wantUser  *User
		wantErr   bool
	}{
		{
			name:     "user found",
			username: "alice",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(int64(1), "alice", "secretpass")
				mock.ExpectQuery(`SELECT\s+id,\s+username,\s+password\s+FROM users\s+WHERE username = \$1`).
					WithArgs("alice").
					WillReturnRows(rows)
			},
			wantUser: &User{ID: 1, Username: "alice", Password: "secretpass"},
			wantErr:  false,
		},
		{
			name:     "user not found",
			username: "missinguser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT\s+id,\s+username,\s+password\s+FROM users\s+WHERE username = \$1`).
					WithArgs("missinguser").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			store := &DBUserStore{DB: db}
			user, err := store.GetUser(tt.username)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUser, user)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBUserStore_UsernameExists(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		mockSetup func(mock sqlmock.Sqlmock)
		want      bool
		wantErr   bool
	}{
		{
			name:     "username exists",
			username: "bob",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM users WHERE username = \$1\)`).
					WithArgs("bob").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "username does not exist",
			username: "charlie",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM users WHERE username = \$1\)`).
					WithArgs("charlie").
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "query error",
			username: "erroruser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM users WHERE username = \$1\)`).
					WithArgs("erroruser").
					WillReturnError(errors.New("query failed"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			store := &DBUserStore{DB: db}
			got, err := store.UsernameExists(tt.username)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBUserStore_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		user      *User
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful create",
			user: &User{Username: "dave", Password: "password123"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(username, password\) VALUES \(\$1, \$2\)`).
					WithArgs("dave", "password123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "create error",
			user: &User{Username: "dave", Password: "password123"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(username, password\) VALUES \(\$1, \$2\)`).
					WithArgs("dave", "password123").
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
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			store := &DBUserStore{DB: db}
			err = store.CreateUser(tt.user)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
