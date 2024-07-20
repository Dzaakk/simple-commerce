package psql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Postgres() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic("Failed to connect to Postgres")
	}

	err = db.Ping()
	if err != nil {
		panic("Error pinging Postgres")
	}

	fmt.Println("Success connect to Postgres")

	return db
}

func Redis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	fmt.Printf("CONNECTION LINK = %v\n", redisAddr)
	fmt.Printf("PASS = %v\n", redisPassword)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis")
	}

	fmt.Println("Success connect to Redis")
	return client
}
