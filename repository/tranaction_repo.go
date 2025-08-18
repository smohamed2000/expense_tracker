package repository

import (
	"tracker/models"

	"gorm.io/gorm"
)

type TransactionRepo struct{ DB *gorm.DB }

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
	GetTotalIncome(userID uint) (float64, error)
	GetTotalExpense(userID uint) (float64, error)
	GetTotalBalance(userID uint) (float64, error)
}

// CreateTransaction saves a new transaction
func (r *TransactionRepo) CreateTransaction(tx *models.Transaction) error {
	return r.DB.Create(tx).Error
}

// GetTransactionsByUserID fetches all transactions for a user
func (r *TransactionRepo) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.DB.Where("user_id = ?", userID).Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTotalIncome returns the total income for a user
func (r *TransactionRepo) GetTotalIncome(userID uint) (float64, error) {
	var total float64
	err := r.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalExpense returns the total expense for a user
func (r *TransactionRepo) GetTotalExpense(userID uint) (float64, error) {
	var total float64
	err := r.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalBalance calculates total balance (income - expense) for a user
func (r *TransactionRepo) GetTotalBalance(userID uint) (float64, error) {
	var total float64
	err := r.DB.Model(&models.Transaction{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0)").
		Scan(&total).Error
	return total, err
}
