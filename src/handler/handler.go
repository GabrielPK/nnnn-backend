package handler

import (
	"encoding/json"
	"net/http"
	"fmt"
	// "context"
	// "time"
	"gorm.io/gorm"
	"nnnn/main/models"
	"golang.org/x/crypto/scrypt"
	"log"
	"encoding/base64"
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
			"password": signUpRequest.Password,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
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

func HashPassword(password string) (string, error) {
    // Generate a random salt
    salt := make([]byte, 16)

    // Generate the hashed password
    dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
    if err != nil {
		log.Fatal(err)
	}

    // Return the salt and the hashed password, encoded in base64 and concatenated
    return base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(dk), nil
}