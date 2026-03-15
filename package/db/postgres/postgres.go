package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func Init(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	for range 5 {
		db, err := sql.Open("postgres", dsn)
		if err == nil && db.Ping() == nil {
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
			db.SetConnMaxIdleTime(1 * time.Minute)

			log.Print("Success connect to Postgres")
			return db, nil
		}

		if db != nil {
			db.Close()
		}

		log.Print("Postgres is not ready, retrying...")
		time.Sleep(5 * time.Second)
	}

	return nil, errors.New("failed to connect to Postgres after multiple attempts")
}
