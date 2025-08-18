package service

import (
	"tracker/models"
	"tracker/repository"
)

type TransactionService struct {
	Repo repository.TransactionRepository
}

// CreateTransaction creates a new transaction
func (t *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	return t.Repo.CreateTransaction(transaction)
}

// GetTransactionsByUserID fetches all transactions for a user
func (t *TransactionService) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	return t.Repo.GetTransactionsByUserID(userID)
}

// GetTotalIncome returns total income for a user
func (t *TransactionService) GetTotalIncome(userID uint) (float64, error) {
	return t.Repo.GetTotalIncome(userID)
}

// GetTotalExpense returns total expense for a user
func (t *TransactionService) GetTotalExpense(userID uint) (float64, error) {
	return t.Repo.GetTotalExpense(userID)
}

// GetTotalBalance returns the total balance for a user
func (t *TransactionService) GetTotalBalance(userID uint) (float64, error) {
	return t.Repo.GetTotalBalance(userID)
}
