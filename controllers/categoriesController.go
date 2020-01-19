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

var categoryData DAL.CategoryData

// NewCategoriesController constructor
func NewCategoriesController(category DAL.CategoryData) {
	categoryData = category
}

// GetCategory /categories/{id}
func GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := categoryData.GetDataByID(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// AddCategory /categories
func AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var category models.Category
	err := decoder.Decode(&category)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := categoryData.AddCategory(&category)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID primitive.ObjectID
	}{id})
}
