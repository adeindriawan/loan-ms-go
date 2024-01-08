package usecase

import (
	"github.com/go-redis/redis"
)

type UseCase struct {
	Cache *redis.Client
}

func NewUseCase(cache *redis.Client) *UseCase {
	return &UseCase{
		Cache: cache,
	}
}