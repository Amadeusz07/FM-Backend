package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there - this page was served using Go \\o/")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got Request Ping Pong")
	w.Header().Set("Content-Type", "application/json")
	result := "Pong"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
