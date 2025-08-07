package model

type UserService interface {
	GetUser(username string) (*User, error)
	UsernameExists(username string) (bool, error)
	CreateUser(user *User) error
}
