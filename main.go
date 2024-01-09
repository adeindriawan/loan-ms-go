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

	cfg := &config.BootstrapConfig{
		DB: db,
		Cache: redisClient,
	}

	appConfig := config.Bootstrap(cfg)

	// Set up your routes with handlers
	router.HandleFunc("/", handlers.HomeHandler)
	router.HandleFunc("/add", handlers.AddUserHandler(appConfig.UserUseCase)).Methods("POST")
	router.HandleFunc("/users", handlers.GetUsersHandler(appConfig.UserUseCase))
	router.HandleFunc("/user/{id}", handlers.GetUserByIDHandler(appConfig.UserUseCase))
	router.HandleFunc("/update", handlers.UpdateUserHandler(appConfig.UserUseCase)).Methods("POST")

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
