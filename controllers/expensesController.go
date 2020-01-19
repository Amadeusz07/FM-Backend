package controllers

import (
	"fmt"
	"net/http"

	"github.com/Amadeusz07/FM-Backend/DAL"
	"github.com/gorilla/mux"
)

var expenseData ExpenseData

// NewExpensesController constructor
func NewExpensesController(db DAL.Database) {
	expenseData = db.NewExpenseData()
}

// GetExpense expenses/{id}
func GetExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Received request for expenses id: %s", vars["id"])
}
