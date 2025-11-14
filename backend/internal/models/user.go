package models

// User represents a user in the system with a hashed password.
type User struct {
	ID       int64
	Username string
	Password string // hashed password
}

// Credentials represents user login/signup data in JSON.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordRequest represents the data needed to change a user's password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
