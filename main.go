package main

import (
	db "Dzaakk/simple-commerce/package/db"
	"fmt"

	customer "Dzaakk/simple-commerce/internal/customer/injector"
	product "Dzaakk/simple-commerce/internal/product/injector"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/injector"
	transaction "Dzaakk/simple-commerce/internal/transaction/injector"

	"github.com/gin-gonic/gin"
)

func main() {
	postgres := db.Postgres()
	redis := db.Redis()
	r := gin.Default()
	fmt.Println("START")

	customer.InitializedService(postgres).Route(&r.RouterGroup, redis)
	product.InitializedService(postgres).Route(&r.RouterGroup, redis)
	shoppingCart.InitializedService(postgres).Route(&r.RouterGroup, redis)
	transaction.InitializedService(postgres).Route(&r.RouterGroup, redis)
	r.Run()
}
