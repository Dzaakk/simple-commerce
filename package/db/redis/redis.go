package redis

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Builder struct {
	cfg Config
}

func NewBuilder() *Builder {
	return &Builder{}
}

func Init() (*redis.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return NewBuilder().
		WithHost(os.Getenv("REDIS_HOST")).
		WithPort(os.Getenv("REDIS_PORT")).
		WithPassword(os.Getenv("REDIS_PASSWORD")).
		WithDB(0).
		Connect()
}

func (b *Builder) WithConfig(cfg Config) *Builder {
	b.cfg = cfg
	return b
}

func (b *Builder) WithHost(host string) *Builder {
	b.cfg.Host = host
	return b
}

func (b *Builder) WithPort(port string) *Builder {
	b.cfg.Port = port
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.cfg.Password = password
	return b
}

func (b *Builder) WithDB(db int) *Builder {
	b.cfg.DB = db
	return b
}

func (b *Builder) Connect() (*redis.Client, error) {
	redisAddr := b.cfg.Host + ":" + b.cfg.Port

	for range 5 {
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: b.cfg.Password,
			DB:       b.cfg.DB,
		})

		if _, err := client.Ping(context.Background()).Result(); err == nil {
			log.Print("Success connect to Redis")
			return client, nil
		}

		client.Close()
		log.Print("Redis is not ready, retrying...")
		time.Sleep(5 * time.Second)
	}

	return nil, errors.New("failed to connect to Redis after multiple attempts")
}
