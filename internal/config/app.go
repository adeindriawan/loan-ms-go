package config

import (
	"database/sql"
	"github.com/go-redis/redis"
	"loan-ms-go/internal/repository"
	"loan-ms-go/internal/usecase"
)

type BootstrapConfig struct {
	DB *sql.DB
	Cache *redis.Client
}

type AppConfig struct {
	UserUseCase *usecase.UserUseCase
}

func Bootstrap(config *BootstrapConfig) *AppConfig {
	userRepository := repository.NewUserRepository(config.DB)
	userUseCase := usecase.NewUserUseCase(config.Cache, userRepository)

	return &AppConfig{
		UserUseCase: userUseCase,
	}
}