/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"myapp/internal"
	"os"
	"log"
	"net/http"
	"myapp/api"
	"github.com/joho/godotenv"
	"fmt"
	"encoding/json"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	internal.InitDB(connStr)

	http.HandleFunc("/register", api.RegisterHandler)
	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/users", api.ListUsersHandler)
	http.HandleFunc("/forgot-password", api.ForgotPasswordHandler)
	http.HandleFunc("/reset-password", api.ResetPasswordHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"message": "OK",
		})
	})

	log.Println("Starting server on :8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
