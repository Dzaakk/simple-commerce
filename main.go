package main

import (
	db "Dzaakk/simple-commerce/package/db"
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
	postgres, err := db.Postgres()
	if err != nil {
		log.Fatalf("error connect to database : %v", err)
	}
	redis := db.Redis()
	r := gin.Default()

	auth.InitializedService(postgres).Route(&r.RouterGroup, redis)
	customer.InitializedService(postgres).Route(&r.RouterGroup, redis)
	product.InitializedService(postgres).Route(&r.RouterGroup, redis)
	seller.InitializedService(postgres).Route(&r.RouterGroup, redis)
	shoppingCart.InitializedService(postgres).Route(&r.RouterGroup, redis)
	transaction.InitializedService(postgres).Route(&r.RouterGroup, redis)
	r.Run()
}
