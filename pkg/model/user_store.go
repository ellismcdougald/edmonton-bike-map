package model

import "database/sql"

// DBUserStore provides methods to interact with the users table in the database.
type DBUserStore struct {
	DB *sql.DB
}

// GetUser retrieves a user by username.
// Returns a pointer to User and an error if the user is not found or query fails.
func (s *DBUserStore) GetUser(username string) (*User, error) {
	query := `
		SELECT
			id,
			username,
			password
		FROM users
		WHERE username = $1
	`

	var user User
	err := s.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UsernameExists checks if a username exists in the database.
// Returns true if the username exists, false otherwise, along with any error encountered.
func (s *DBUserStore) UsernameExists(username string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := s.DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// CreateUser inserts a new user into the database.
// Returns any error encountered during the insert.
func (s *DBUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	_, err := s.DB.Exec(query, user.Username, user.Password)
	return err
}
