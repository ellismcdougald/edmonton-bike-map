package model

// UserService defines the interface for managing users in a data store.
// It supports getting an individual user, checking whether a username exists, and creating a user.
type UserService interface {
	GetUser(username string) (*User, error)
	UsernameExists(username string) (bool, error)
	CreateUser(user *User) error
}
