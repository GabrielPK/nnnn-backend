package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"nnnn/main/models"

	"golang.org/x/crypto/scrypt"
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
		hashedPassword, err = hashPassword(signUpRequest.Password)
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
		if err := comparePasswords(user.Password, logInRequest.Password); err != nil {
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

func hashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Generate the hashed password
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	// Return the salt and the hashed password, encoded in base64 and concatenated
	return base64.StdEncoding.EncodeToString(salt) + base64.StdEncoding.EncodeToString(dk), nil
}

func comparePasswords(hashedPwd, plainPwd string) error {
	// Decode the salt (first 24 characters after base64 encoding of 16 bytes)
	salt, err := base64.StdEncoding.DecodeString(hashedPwd[:24])
	if err != nil {
		return err
	}

	// Decode the stored hash (the rest of the string)
	storedHash, err := base64.StdEncoding.DecodeString(hashedPwd[24:])
	if err != nil {
		return err
	}

	// Hash the provided password using the same salt
	hash, err := scrypt.Key([]byte(plainPwd), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return err
	}

	// Compare the hashes
	if !bytes.Equal(hash, storedHash) {
		return fmt.Errorf("password does not match")
	}
	return nil
}
