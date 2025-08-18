package handler

import (
	"encoding/json"

	"net/http"
	"strconv"
	"tracker/middleware"
	"tracker/models"
	"tracker/service"
)

type BudgetHandler struct {
	Service *service.BudgetService
}

// CreateBudget creates a budget for the logged-in user
func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	var budget models.Budget
	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	budget.UserID = userID

	if err := h.Service.CreateBudget(&budget); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budget)
}

// GetBudgetsByUserID fetches budgets for the logged-in user
func (h *BudgetHandler) GetBudgetsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	budgets, err := h.Service.GetBudgetsByUserID(userID)
	if err != nil {
		http.Error(w, "failed to fetch budgets", http.StatusInternalServerError)
		return
	}

	if len(budgets) == 0 {
		http.Error(w, "no budgets found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

// UpdateBudget updates a budget for the logged-in user
func (h *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	var budget models.Budget
	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id query parameter is missing", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid budget ID", http.StatusBadRequest)
		return
	}
	budget.ID = uint(idInt)

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	budget.UserID = userID

	if err := h.Service.UpdateBudget(&budget); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budget)
}

// DeleteBudget deletes a budget for the logged-in user
func (h *BudgetHandler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id query parameter is missing", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid budget ID", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.Service.DeleteBudget(uint(idInt), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
