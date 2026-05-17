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

	postgres, err := postgres.Init(postgres.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}

	redis, err := redis.Init()
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

	r := gin.Default()
	r.Use(requestid.RequestID())
	r.Use(metricsMiddleware.HTTPMiddleware())
	r.Use(logMiddleware.RequestLogger(logging.NewLokiClientFromEnv()))
	r.Use(middleware.ErrorHandler())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	health.NewHandler(postgres, redis).Route(r)

	auth.InitializedService(postgres, redis, rabbitClient).Route(&r.RouterGroup)
	user.InitializedService(postgres).Route(&r.RouterGroup)
	catalog.InitializedService(postgres).Route(&r.RouterGroup)
	cart.InitializedService(postgres).Route(&r.RouterGroup)
	order.InitializedService(postgres).Route(&r.RouterGroup)
	transaction.InitializedService(postgres).Route(&r.RouterGroup)

	r.Run()
}
