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

// GetCategories /categories
func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	result := categoryData.GetCategories(projectId)
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// GetCategory /categories/:id
func GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := categoryData.GetDataByID(projectId, id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// AddCategory /categories
func AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	decoder := json.NewDecoder(r.Body)
	var category models.Category
	err := decoder.Decode(&category)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := categoryData.AddCategory(userId, projectId, &category)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID primitive.ObjectID
	}{id})
}

// UpdateCategory /categories/:id
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	decoder := json.NewDecoder(r.Body)
	var category models.Category
	if err := decoder.Decode(&category); err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	categoryData.UpdateCategory(projectId, id, &category)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteCategory /categories/:id
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isAny := expenseData.IsAnyInCategory(projectId, id)
	if isAny {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	categoryData.DeleteCategory(projectId, id)
	w.WriteHeader(http.StatusAccepted)
}

func GetCategoriesSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId, _ := primitive.ObjectIDFromHex(r.Header.Get("selectedProjectId"))
	result := expenseData.GetSummary(projectId)
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
