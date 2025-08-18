package repository

import (
	"errors"

	"tracker/models"

	"gorm.io/gorm"
)

type UserRepo struct{ DB *gorm.DB }

// GetUserByEmail fetches a user by email
func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // let service distinguish "not found"
		}
		return nil, err // some other DB error
	}
	return &user, nil
}

// CreateUser inserts a new user
func (r *UserRepo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}
