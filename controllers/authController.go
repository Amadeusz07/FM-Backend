package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../DAL"
	"../dtos"
	"../models"
	authService "../services/auth"

	"golang.org/x/crypto/bcrypt"
)

var userData DAL.UserData

func NewAuthController(user DAL.UserData) {
	userData = user
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
	response := dtos.LoginResponse{token, expDate}
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

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
