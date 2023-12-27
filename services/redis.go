package services

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password
		DB:       0,  // Default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Redis:", pong)

	return client
}
