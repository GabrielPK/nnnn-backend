package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nnnn/main/models"

	"gorm.io/gorm"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var hashedPassword string
		hashedPassword, err = HashPassword(signUpRequest.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if username exists
		var user models.User
		db.Where("username = ?", signUpRequest.Username).First(&user)

		if user.Username != "" { // username exists
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}

		db.Create(&models.User{Username: signUpRequest.Username, Password: hashedPassword})

		// For example, just send them back as a response
		responseData := map[string]string{
			"username": signUpRequest.Username,
			"password": hashedPassword,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LogInHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Retrieve the user from the database
		var user models.User
		result := db.Where("username = ?", logInRequest.Username).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Verify the password
		if err := ComparePasswords(user.Password, logInRequest.Password); err != nil {
			// Incorrect password or hashing error
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		// Login successful
		responseData := map[string]string{
			"message": "Login successful",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}
