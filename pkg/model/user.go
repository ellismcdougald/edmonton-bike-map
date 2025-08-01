package model

import "database/sql"

type User struct {
	ID       int
	Username string
	Password string // hashed
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsernameExists checks whether a given username already exists in the database
func UsernameExists(db *sql.DB, username string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (user *User) Create(db *sql.DB) error {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	_, err := db.Exec(query, user.Username, user.Password)
	return err
}

func GetUser(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT
			id,
			username,
			password
		FROM users
		WHERE username = $1
	`

	var user User
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}