package controllers

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"

	"../DAL"
	"../dtos"
	"../models"
	authService "../services/auth"

	"golang.org/x/crypto/bcrypt"
)

var userData DAL.UserData
var projectDataAuth DAL.ProjectData

func NewAuthController(user DAL.UserData, project DAL.ProjectData) {
	userData = user
	projectDataAuth = project
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var creds models.User
	if err := decoder.Decode(&creds); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := userData.GetUserByEmail(creds.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid password")
		return
	}
	expDate := time.Now().Add(time.Minute * 30)
	token, err := authService.GenerateJWT(user.ID, user.Email, expDate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	response := dtos.LoginResponse{token, expDate, user.Email, ""}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user models.User
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Password = string(hash)
	err = userData.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SelectProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	decoder := json.NewDecoder(r.Body)
	var request dtos.ChangeProjectRequest
	if err := decoder.Decode(&request); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := userData.GetUserById(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !projectDataAuth.IsUserAssignedToProject(request.ProjectId, userId) && !projectDataAuth.IsOwnerOfProject(request.ProjectId, userId) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expDate := time.Now().Add(time.Minute * 30)
	token, err := authService.GenerateJWTWithProject(user.ID, user.Email, expDate, request.ProjectId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	response := dtos.LoginResponse{token, expDate, user.Email, request.ProjectId.Hex()}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
