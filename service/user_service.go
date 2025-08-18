package service

import (
	"errors"

	"tracker/middleware"
	"tracker/models"
	"tracker/utils"

	"gorm.io/gorm"
)

// Define the interface here so we don't depend on repository package types.
type UserRepo interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService struct {
	Repo UserRepo
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(user *models.User) error {
	if user == nil {
		return errors.New("nil user")
	}

	// Check if a user with this email already exists.
	_, err := s.Repo.GetUserByEmail(user.Email)
	switch {
	case err == nil:
		// Found existing user -> duplicate
		return ErrEmailAlreadyInUse
	case errors.Is(err, gorm.ErrRecordNotFound):
		// Not found -> continue to create
	default:
		// Some other DB error
		return err
	}

	// Hash password
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hash

	// Save user
	if err := s.Repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}

// LoginUser validates credentials and generates a JWT
func (s *UserService) LoginUser(req models.LoginRequest) (string, error) {
	// Look up user
	u, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		// Don't leak which part failed
		return "", ErrInvalidCredentials
	}

	// Check password
	if err := utils.ComparePassword(u.Password, req.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Issue token
	token, err := middleware.GenerateJWT(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
