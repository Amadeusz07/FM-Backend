package controllers

import (
	"../DAL"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var projectData DAL.ProjectData

func NewProjectsController(projects DAL.ProjectData) {
	projectData = projects
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	decoder := json.NewDecoder(r.Body)
	var project models.Project
	err := decoder.Decode(&project)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := projectData.AddProject(ownerId, &project)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID primitive.ObjectID
	}{id})
}

func GetProjectByOwnerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))

	projects := projectData.GetProjectsForOwner(userId)
	assignedProjects := projectData.GetAssignedProjects(userId)

	allProjects := append(projects, assignedProjects...)
	if allProjects == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allProjects)
	}
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	decoder := json.NewDecoder(r.Body)
	var project models.Project
	if err := decoder.Decode(&project); err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectData.UpdateProject(ownerId, id, &project)
	w.WriteHeader(http.StatusNoContent)
}

func AssignUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var projectUser models.ProjectUser
	if err := decoder.Decode(&projectUser); err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectData.AssignUser(id, projectUser.UserId)

	w.WriteHeader(http.StatusNoContent)
}

func UnAssignUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var projectUser models.ProjectUser
	if err := decoder.Decode(&projectUser); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectData.UnAssignUser(id, projectUser.UserId)

	w.WriteHeader(http.StatusNoContent)
}
