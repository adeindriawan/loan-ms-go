package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"loan-ms-go/internal/handlers"
	"loan-ms-go/internal/config"
	"loan-ms-go/services"
)

func main() {
	router := mux.NewRouter()

	// Initialize database connections
	db := services.InitMySQL()
	redisClient := services.InitRedis()

	config.Bootstrap(&config.BootstrapConfig{
		DB: db,
	})

	// Set up your routes with handlers
	router.HandleFunc("/", handlers.HomeHandler)
	router.HandleFunc("/add", handlers.AddUserHandler(db, redisClient)).Methods("POST")
	router.HandleFunc("/users", handlers.GetUsersHandler(db))
	router.HandleFunc("/user/{id}", handlers.GetUserHandler(db, redisClient))
	router.HandleFunc("/update", handlers.UpdateUserHandler(db, redisClient)).Methods("POST")

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on %s...\n", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
