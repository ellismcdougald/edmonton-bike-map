package sqlrepo

import (
	"database/sql"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// SQLUserRepository implements UserRepository using a SQL database.
type SQLUserRepository struct {
	DB *sql.DB
}

// NewSQLUserRepository creates a new SQLUserRepository.
func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{DB: db}
}

// GetByUsername retrieves a user by their username.
func (s *SQLUserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
			SELECT
				id,
				username,
				password
			FROM users
			WHERE username = $1
		`

	var user models.User
	err := s.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by their ID.
func (s *SQLUserRepository) GetByID(id int64) (*models.User, error) {
	query := `
			SELECT
				id,
				username,
				password
			FROM users
			WHERE id = $1
		`

	var user models.User
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser inserts a new user into the database.
// Returns any error encountered during the insert.
func (s *SQLUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	_, err := s.DB.Exec(query, user.Username, user.Password)
	return err
}

// UsernameExists checks if a username exists in the database.
// Returns true if the username exists, false otherwise, along with any error encountered.
func (s *SQLUserRepository) UsernameExists(username string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := s.DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// UpdatePassword updates a user's password in the database.
func (s *SQLUserRepository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `
		UPDATE users SET password = $1 WHERE id = $2
	`
	_, err := s.DB.Exec(query, hashedPassword, userID)
	return err
}
