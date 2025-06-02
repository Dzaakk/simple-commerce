package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func Redis() *redis.Client {
	redisOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
		redisPassword := os.Getenv("REDIS_PASSWORD")

		client := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})

		_, err = client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Print("Success connect to Redis")
		redisClient = client
	})

	return redisClient
}
