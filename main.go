package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"loan-ms-go/internal/config"
	"loan-ms-go/services"
)

func main() {
	router := mux.NewRouter()

	db := services.InitMySQL()
	redisClient := services.InitRedis()

	cfg := &config.BootstrapConfig{
		DB: db,
		Cache: redisClient,
		Router: router,
	}
	config.Bootstrap(cfg)

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
