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

type Builder struct {
	cfg Config
}

func NewBuilder() *Builder {
	return &Builder{}
}

func Init(cfg Config) (*sql.DB, error) {
	return NewBuilder().
		WithConfig(cfg).
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

func (b *Builder) WithDBName(dbName string) *Builder {
	b.cfg.DBName = dbName
	return b
}

func (b *Builder) WithUser(user string) *Builder {
	b.cfg.User = user
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.cfg.Password = password
	return b
}

func (b *Builder) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		b.cfg.Host, b.cfg.Port, b.cfg.User, b.cfg.Password, b.cfg.DBName,
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
