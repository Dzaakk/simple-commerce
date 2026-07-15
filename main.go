package main

import (
	emailQueue "Dzaakk/simple-commerce/internal/email/queue"
	emailService "Dzaakk/simple-commerce/internal/email/service"
	postgres "Dzaakk/simple-commerce/package/db/postgres"
	redis "Dzaakk/simple-commerce/package/db/redis"
	"log"
	"os"

	auth "Dzaakk/simple-commerce/internal/auth/route"
	cart "Dzaakk/simple-commerce/internal/cart/route"
	catalog "Dzaakk/simple-commerce/internal/catalog/route"
	"Dzaakk/simple-commerce/internal/health"
	"Dzaakk/simple-commerce/internal/middleware"
	logMiddleware "Dzaakk/simple-commerce/internal/middleware/logging"
	metricsMiddleware "Dzaakk/simple-commerce/internal/middleware/metrics"
	requestid "Dzaakk/simple-commerce/internal/middleware/requestid"
	order "Dzaakk/simple-commerce/internal/order/route"
	transaction "Dzaakk/simple-commerce/internal/transaction/route"
	user "Dzaakk/simple-commerce/internal/user/route"
	"Dzaakk/simple-commerce/package/logging"
	"Dzaakk/simple-commerce/package/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	postgresDB, err := postgres.NewBuilder().
		WithHost(os.Getenv("POSTGRES_HOST")).
		WithPort(os.Getenv("POSTGRES_PORT")).
		WithDBName(os.Getenv("POSTGRES_DB")).
		WithUser(os.Getenv("POSTGRES_USER")).
		WithPassword(os.Getenv("POSTGRES_PASSWORD")).
		Connect()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := redis.NewBuilder().
		WithHost(os.Getenv("REDIS_HOST")).
		WithPort(os.Getenv("REDIS_PORT")).
		WithPassword(os.Getenv("REDIS_PASSWORD")).
		WithDB(0).
		Connect()
	if err != nil {
		log.Fatalf("error connect to redis : %v", err)
	}

	var rabbitClient *rabbitmq.Client
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL != "" {
		rabbitClient, err = rabbitmq.Init(rabbitURL)
		if err != nil {
			log.Printf("rabbitmq queue disabled: %v", err)
		} else {
			defer rabbitClient.Close()

			if err := emailQueue.StartActivationEmailConsumer(rabbitClient, emailService.NewEmailService()); err != nil {
				log.Printf("failed to start activation email consumer: %v", err)
			}
		}
	}

	r := gin.New()
	r.Use(requestid.RequestID())
	r.Use(metricsMiddleware.HTTPMiddleware())
	r.Use(logMiddleware.RequestLogger(logging.NewLogger("http", "api")))
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	health.NewHandler(postgresDB, redisClient).Route(r)

	auth.InitializedService(postgresDB, redisClient, rabbitClient).Route(&r.RouterGroup)
	user.InitializedService(postgresDB).Route(&r.RouterGroup)
	catalog.InitializedService(postgresDB, redisClient).Route(&r.RouterGroup)
	cart.InitializedService(postgresDB).Route(&r.RouterGroup)
	order.InitializedService(postgresDB).Route(&r.RouterGroup)
	transaction.InitializedService(postgresDB).Route(&r.RouterGroup)

	r.Run()
}
