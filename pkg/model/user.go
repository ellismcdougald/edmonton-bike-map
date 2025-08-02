package model

import "database/sql"

// User represents a user in the system with a hashed password.
type User struct {
	ID       int
	Username string
	Password string // hashed password
}

// Credentials represents user login/signup data in JSON.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsernameExists checks whether a given username already exists in the database.
// Returns true if the username exists, false otherwise, along with any error encountered.
func UsernameExists(db *sql.DB, username string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// Create inserts a new user into the database.
func (user *User) Create(db *sql.DB) error {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	_, err := db.Exec(query, user.Username, user.Password)
	return err
}

// GetUser retrieves a user by username from the database.
// Returns a pointer to the User or an error if the user does not exist or query fails.
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
