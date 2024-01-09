package config

import (
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"loan-ms-go/internal/repository"
	"loan-ms-go/internal/usecase"
	"loan-ms-go/internal/handlers"
)

type BootstrapConfig struct {
	DB *sql.DB
	Cache *redis.Client
	Router *mux.Router
}

func Bootstrap(config *BootstrapConfig) {
	config.Router.HandleFunc("/", handlers.HomeHandler)

	userRepository := repository.NewUserRepository(config.DB)
	userUseCase := usecase.NewUserUseCase(config.Cache, userRepository)

	config.Router.HandleFunc("/users/{id}/details", handlers.GetUserByIDHandler(userUseCase))
	config.Router.HandleFunc("/users/add", handlers.AddUserHandler(userUseCase)).Methods("POST")
	config.Router.HandleFunc("/users", handlers.GetUsersHandler(userUseCase))
	config.Router.HandleFunc("/users/update", handlers.UpdateUserHandler(userUseCase)).Methods("POST")
}