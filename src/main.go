package main

import (
	"log"
	"net/http"

	"fmt"
	"nnnn/main/handler"
	"time"

	"github.com/gorilla/mux"

	"nnnn/main/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db := database.InitializeDB()

	// Set up the router.
	r := mux.NewRouter()

	// Set up routes.
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/signup", handler.SignUpHandler(db))
	r.HandleFunc("/login", handler.LogInHandler(db))
	r.HandleFunc("/send-notification", handler.SendNotificationHandler(db))
	r.HandleFunc("/list-users", handler.ListUsersHandler(db))
	r.HandleFunc("/get-notifications-for-user", handler.GetNotificationsForUserHandler(db))

	// Initialize the server with some basic configurations.
	srv := &http.Server{
		Handler:      r,                // Use the mux router as the handler
		Addr:         "127.0.0.1:8080", // Bind address and port for the server
		WriteTimeout: 15 * time.Second, // Max duration for writing responses
		ReadTimeout:  15 * time.Second, // Max duration for reading request bodies
		IdleTimeout:  60 * time.Second, // Max duration for idle connections
	}

	// Start the server.
	fmt.Println("Starting server on http://127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
