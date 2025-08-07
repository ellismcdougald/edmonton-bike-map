package model

import "database/sql"

type DBUserStore struct {
	DB *sql.DB
}

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

func (s *DBUserStore) UsernameExists(username string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := s.DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (s *DBUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	_, err := s.DB.Exec(query, user.Username, user.Password)
	return err
}
