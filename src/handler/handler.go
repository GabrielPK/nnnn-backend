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
		hashedPassword, err = SaltAndHashPassword(signUpRequest.Password)
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
			"password": "uhoh",
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

type SendNotificationRequest struct {
	SenderId   uint   `json:"sender_id"`
	ReceiverId uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func SendNotificationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var sendNotificationRequest SendNotificationRequest

		// Decode the JSON body into the struct
		err := json.NewDecoder(r.Body).Decode(&sendNotificationRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if sender_id exists
		var sender models.User
		db.Where("id = ?", sendNotificationRequest.SenderId).First(&sender)

		if sender.Username == "" { // sender_id does not exist
			http.Error(w, "Sender does not exist", http.StatusBadRequest)
			return
		}

		// check if receiver_id exists
		var receiver models.User
		db.Where("id = ?", sendNotificationRequest.ReceiverId).First(&receiver)

		if receiver.Username == "" { // receiver_id does not exist
			http.Error(w, "Receiver does not exist", http.StatusBadRequest)
			return
		}

		db.Create(&models.Notification{SenderId: sendNotificationRequest.SenderId, ReceiverId: sendNotificationRequest.ReceiverId, Content: sendNotificationRequest.Content})

		// For example, just send them back as a response
		responseData := map[string]string{
			"sender_id":   fmt.Sprint(sendNotificationRequest.SenderId),
			"receiver_id": fmt.Sprint(sendNotificationRequest.ReceiverId),
			"content":     sendNotificationRequest.Content,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

func ListUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		db.Find(&users)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

type GetNotificationsForUserRequest struct {
	UserId uint `json:"user_id"`
}

func GetNotificationsForUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getNotificationsForUserRequest GetNotificationsForUserRequest

		// Decode the JSON body into the struct
		err := json.NewDecoder(r.Body).Decode(&getNotificationsForUserRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var notifications []models.Notification
		db.Where("receiver_id = ?", getNotificationsForUserRequest.UserId).Find(&notifications)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}
}
