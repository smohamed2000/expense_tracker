package handler

import (
	"encoding/json"
	"net/http"
	"tracker/middleware"
	"tracker/models"
	"tracker/service"
)

type TransactionHandler struct {
	Service *service.TransactionService
}

// CreateTransaction creates a transaction for the logged-in user
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "could not get user id", http.StatusUnauthorized)
		return
	}

	transaction.UserID = userID

	if err := h.Service.CreateTransaction(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// GetTransactionsByUserID returns transactions for the logged-in user
func (h *TransactionHandler) GetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "could not get user id", http.StatusUnauthorized)
		return
	}

	transactions, err := h.Service.GetTransactionsByUserID(userID)
	if err != nil {
		http.Error(w, "failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	if len(transactions) == 0 {
		http.Error(w, "no transactions found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetTotalBalance returns the total balance for the logged-in user
func (h *TransactionHandler) GetTotalBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "could not get user id", http.StatusUnauthorized)
		return
	}

	totalBalance, err := h.Service.GetTotalBalance(userID)
	if err != nil {
		http.Error(w, "could not calculate total balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": totalBalance})
}
