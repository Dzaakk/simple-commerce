package redis

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	err         error
)

func Init() (*redis.Client, error) {
	redisOnce.Do(func() {
		err = godotenv.Load()
		if err != nil {
			return
		}

		redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
		redisPassword := os.Getenv("REDIS_PASSWORD")

		for range 5 {
			client := redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: redisPassword,
				DB:       0,
			})

			_, err = client.Ping(context.Background()).Result()
			if err == nil {
				log.Print("Success connect to Redis")
				redisClient = client
				return
			}

			client.Close()
			log.Print("Redis is not ready, retrying...")
			time.Sleep(5 * time.Second)
		}

		err = errors.New("failed to connect to Redis after multiple attempts")
	})

	return redisClient, err
}
