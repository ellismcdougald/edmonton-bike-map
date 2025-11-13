package service

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepo,
	}
}

// GetUserByUsername retrieves a user by their username using the UserRepository.
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepository.GetByUsername(username)
}

// CreateUser creates a new user using the UserRepository.
func (s *UserService) CreateUser(user *models.User) error {
	return s.UserRepository.Create(user)
}

// UsernameExists checks if a username already exists using the UserRepository.
func (s *UserService) UsernameExists(username string) (bool, error) {
	return s.UserRepository.UsernameExists(username)
}
