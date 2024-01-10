package usecase

import (
	"loan-ms-go/services"
	"github.com/go-redis/redis"
)

type UseCase struct {
	Cache *redis.Client
	Logger *services.Logger
}

func NewUseCase(cache *redis.Client, logger *services.Logger) *UseCase {
	return &UseCase{
		Cache: cache,
		Logger: logger,
	}
}