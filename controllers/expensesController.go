package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../DAL"
	"../models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var expenseData DAL.ExpenseData

// NewExpensesController constructor
func NewExpensesController(expense DAL.ExpenseData) {
	expenseData = expense
}

// GetExpense /expenses/{id}
func GetExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := expenseData.GetDataByID(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// AddExpense /expenses
func AddExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var expense models.Expense
	err := decoder.Decode(&expense)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := expenseData.AddExpense(&expense)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID primitive.ObjectID
	}{id})
}
