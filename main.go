package main

import (
	postgres "Dzaakk/simple-commerce/package/db/postgres"
	redis "Dzaakk/simple-commerce/package/db/redis"
	"log"

	auth "Dzaakk/simple-commerce/internal/auth/injector"
	customer "Dzaakk/simple-commerce/internal/customer/injector"
	product "Dzaakk/simple-commerce/internal/product/injector"
	seller "Dzaakk/simple-commerce/internal/seller/injector"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/injector"
	transaction "Dzaakk/simple-commerce/internal/transaction/injector"

	"github.com/gin-gonic/gin"
)

func main() {
	postgres, err := postgres.Init()
	if err != nil {
		log.Fatalf("error connect to database : %v", err)
	}
	redis, err := redis.Init()
	if err != nil {
		log.Fatalf("error connect to redis : %v", err)
	}

	r := gin.Default()

	auth.InitializedService(postgres, redis).Route(&r.RouterGroup)
	customer.InitializedService(postgres).Route(&r.RouterGroup, redis)
	product.InitializedService(postgres).Route(&r.RouterGroup, redis)
	seller.InitializedService(postgres).Route(&r.RouterGroup, redis)
	shoppingCart.InitializedService(postgres).Route(&r.RouterGroup, redis)
	transaction.InitializedService(postgres).Route(&r.RouterGroup, redis)
	r.Run()
}
