package main

import (
	postgres "Dzaakk/simple-commerce/package/db/postgres"
	redis "Dzaakk/simple-commerce/package/db/redis"
	"log"
	"os"

	auth "Dzaakk/simple-commerce/internal/auth/route"
	cart "Dzaakk/simple-commerce/internal/cart/route"
	catalog "Dzaakk/simple-commerce/internal/catalog/route"
	order "Dzaakk/simple-commerce/internal/order/route"
	transaction "Dzaakk/simple-commerce/internal/transaction/route"
	user "Dzaakk/simple-commerce/internal/user/route"
	"Dzaakk/simple-commerce/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	auth.InitializedService(postgres, redis).Route(&r.RouterGroup)
	user.InitializedService(postgres).Route(&r.RouterGroup)
	catalog.InitializedService(postgres).Route(&r.RouterGroup)
	cart.InitializedService(postgres).Route(&r.RouterGroup)
	order.InitializedService(postgres).Route(&r.RouterGroup)
	transaction.InitializedService(postgres).Route(&r.RouterGroup)

	r.Run()
}
