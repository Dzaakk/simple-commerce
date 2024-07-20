package main

import (
	db "Dzaakk/synapsis/package/db"
	"fmt"

	customer "Dzaakk/synapsis/internal/customer/injector"
	product "Dzaakk/synapsis/internal/product/injector"
	shoppingCart "Dzaakk/synapsis/internal/shopping_cart/injector"
	transaction "Dzaakk/synapsis/internal/transaction/injector"

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
