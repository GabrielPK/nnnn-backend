package main

import (
	"log"
	"net/http"
	// "os"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"nnnn/main/handler"
)

func main() {
	// Set up the router.
	r := mux.NewRouter()

	// Set up routes.
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/signup", handler.SignUpHandler)
	r.HandleFunc("/login", handler.LogInHandler)
	
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
