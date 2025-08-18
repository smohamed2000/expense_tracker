package repository

import (
	"tracker/models"

	"gorm.io/gorm"
)

type BudgetRepo struct{ DB *gorm.DB }

type BudgetRepository interface {
	CheckDuplicateBudget(budget *models.Budget) bool
	CreateBudget(budget *models.Budget) error
	GetBudgetsByUserID(userID uint) ([]models.Budget, error)
	UpdateBudget(budget *models.Budget) error
	CheckBudgetExistsForUser(id uint, userID uint) bool
	DeleteBudget(id uint) error
}

// CheckDuplicateBudget checks if a budget category already exists for the user
func (r *BudgetRepo) CheckDuplicateBudget(budget *models.Budget) bool {
	var count int64
	r.DB.Model(&models.Budget{}).
		Where("category = ? AND user_id = ?", budget.Category, budget.UserID).
		Count(&count)
	return count > 0
}

// CreateBudget inserts a new budget
func (r *BudgetRepo) CreateBudget(budget *models.Budget) error {
	return r.DB.Create(budget).Error
}

// GetBudgetsByUserID fetches all budgets for a user
func (r *BudgetRepo) GetBudgetsByUserID(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	if err := r.DB.Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

// UpdateBudget updates a budget
func (r *BudgetRepo) UpdateBudget(budget *models.Budget) error {
	return r.DB.Save(budget).Error
}

// CheckBudgetExistsForUser checks if a budget exists for a given user
func (r *BudgetRepo) CheckBudgetExistsForUser(id uint, userID uint) bool {
	var count int64
	r.DB.Model(&models.Budget{}).
		Where("id = ? AND user_id = ?", id, userID).
		Count(&count)
	return count > 0
}

// DeleteBudget deletes a budget by ID
func (r *BudgetRepo) DeleteBudget(id uint) error {
	return r.DB.Where("id = ?", id).Delete(&models.Budget{}).Error
}
