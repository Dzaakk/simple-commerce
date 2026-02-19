package main

import (
	postgres "Dzaakk/simple-commerce/package/db/postgres"
	redis "Dzaakk/simple-commerce/package/db/redis"
	"log"

	auth "Dzaakk/simple-commerce/internal/auth/route"
	product "Dzaakk/simple-commerce/internal/product/injector"
	seller "Dzaakk/simple-commerce/internal/seller/injector"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/injector"
	transaction "Dzaakk/simple-commerce/internal/transaction/injector"
	user "Dzaakk/simple-commerce/internal/user/route"

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
	user.InitializedService(postgres, redis).Route(&r.RouterGroup)
	product.InitializedService(postgres, redis).Route(&r.RouterGroup)
	seller.InitializedService(postgres).Route(&r.RouterGroup, redis)
	shoppingCart.InitializedService(postgres, redis).Route(&r.RouterGroup)
	transaction.InitializedService(postgres, redis).Route(&r.RouterGroup)
	r.Run()
}
