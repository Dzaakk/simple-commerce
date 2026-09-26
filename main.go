package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authroute "Dzaakk/simple-commerce/internal/auth/route"
	catalogroute "Dzaakk/simple-commerce/internal/catalog/route"
	"Dzaakk/simple-commerce/internal/health"
	"Dzaakk/simple-commerce/internal/middleware"
	userroute "Dzaakk/simple-commerce/internal/user/route"
	"Dzaakk/simple-commerce/package/db/postgres"
	redisdb "Dzaakk/simple-commerce/package/db/redis"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := postgres.NewBuilder().
		WithHost(os.Getenv("POSTGRES_HOST")).
		WithPort(os.Getenv("POSTGRES_PORT")).
		WithDBName(os.Getenv("POSTGRES_DB")).
		WithUser(os.Getenv("POSTGRES_USER")).
		WithPassword(os.Getenv("POSTGRES_PASSWORD")).
		Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var redisClient *redis.Client
	if os.Getenv("REDIS_HOST") != "" {
		redisClient, err = redisdb.NewBuilder().
			WithHost(os.Getenv("REDIS_HOST")).
			WithPort(os.Getenv("REDIS_PORT")).
			WithPassword(os.Getenv("REDIS_PASSWORD")).
			Connect()
		if err != nil {
			log.Printf("redis unavailable; catalog v2 will fall back to postgres: %v", err)
		}
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	router := newRouter(db, redisClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func newRouter(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.ErrorHandler())
	health.NewHandler(db, redisClient).Route(router)
	authroute.New(db).Route(&router.RouterGroup)
	userroute.Route(&router.RouterGroup, db)
	catalogroute.Route(&router.RouterGroup, db, redisClient)
	return router
}
