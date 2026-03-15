package main

import (
	postgres "Dzaakk/simple-commerce/package/db/postgres"
	redis "Dzaakk/simple-commerce/package/db/redis"
	"log"
	"os"

	auth "Dzaakk/simple-commerce/internal/auth/route"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/injector"
	transaction "Dzaakk/simple-commerce/internal/transaction/injector"
	user "Dzaakk/simple-commerce/internal/user/route"

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

	auth.InitializedService(postgres, redis).Route(&r.RouterGroup)
	user.InitializedService(postgres, redis).Route(&r.RouterGroup)
	shoppingCart.InitializedService(postgres, redis).Route(&r.RouterGroup)
	transaction.InitializedService(postgres, redis).Route(&r.RouterGroup)
	r.Run()
}
