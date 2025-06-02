package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	postgres     *sql.DB
	postgresOnce sync.Once
	err          error
)

func Postgres() (*sql.DB, error) {
	postgresOnce.Do(func() {
		err = godotenv.Load()
		if err != nil {
			return
		}

		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		dbname := os.Getenv("POSTGRES_DB")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")

		connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		for range 5 {
			db, err := sql.Open("postgres", connectionString)
			if err == nil && db.Ping() == nil {
				log.Print("Success connect to Postgres")
				postgres = db
				return
			}
			log.Print("Postgres is not ready, retrying...")
			time.Sleep(5 * time.Second)
		}

		err = errors.New("failed to connect to Postgres after multiple attempts")

	})
	return postgres, err
}
