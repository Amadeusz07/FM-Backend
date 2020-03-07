package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

// GetExpenses /expenses?count=int64
func GetExpenses(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.ParseInt(r.URL.Query().Get("count"), 10, 16)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	result := expenseData.GetLastHistory(userId, count, date)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetExpense /expenses/{id}
func GetExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	result := expenseData.GetDataByID(userId, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// AddExpense /expenses
func AddExpense(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var expense models.Expense
	err := decoder.Decode(&expense)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	id := expenseData.AddExpense(userId, &expense)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID primitive.ObjectID
	}{id})
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expenseData.DeleteExpense(userId, id)
	w.WriteHeader(http.StatusAccepted)
}
