package handler

import (
	"encoding/json"
	"net/http"
	"fmt"
	// "database/sql"
	// "github.com/mattn/go-sqlite3"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var signUpRequest SignUpRequest

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&signUpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For example, just send them back as a response
	responseData := map[string]string{
		"username": signUpRequest.Username,
		"password": signUpRequest.Password,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {	
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var logInRequest LogInRequest

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&logInRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For example, just send them back as a response
	responseData := map[string]string{
		"username": logInRequest.Username,
		"password": logInRequest.Password,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}