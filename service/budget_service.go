package service

import (
	"errors"
	
	"log"
	"tracker/models"
	"tracker/repository"
)

var ErrBudgetNotFound = errors.New("budget not found")

type BudgetService struct {
	Repo repository.BudgetRepository
}

// CreateBudget creates a new budget
func (b *BudgetService) CreateBudget(budget *models.Budget) error {
	return b.Repo.CreateBudget(budget)
}

// GetBudgetsByUserID fetches all budgets for a user
func (b *BudgetService) GetBudgetsByUserID(userID uint) ([]models.Budget, error) {
	budgets, err := b.Repo.GetBudgetsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return budgets, nil
}

// UpdateBudget updates a budget
func (b *BudgetService) UpdateBudget(budget *models.Budget) error {
	return b.Repo.UpdateBudget(budget)
}

// DeleteBudget deletes a budget for a user
func (b *BudgetService) DeleteBudget(id uint, userID uint) error {
	// Ensure budget exists and belongs to user
	ok := b.Repo.CheckBudgetExistsForUser(id, userID)
	if !ok {
		return ErrBudgetNotFound
	}

	log.Printf("Budget with ID %d found for user %d, proceeding to delete", id, userID)
	return b.Repo.DeleteBudget(id)
}
