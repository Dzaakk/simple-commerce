package main

import (
	psql "Dzaakk/synapsis/package/db"
	"fmt"

	customer "Dzaakk/synapsis/internal/customer/injector"
	product "Dzaakk/synapsis/internal/product/injector"

	"github.com/gin-gonic/gin"
)

func main() {
	db := psql.DB()
	r := gin.Default()
	fmt.Println("START")

	customer.InitializedService(db).Route(&r.RouterGroup)
	product.InitializedService(db).Route(&r.RouterGroup)
	r.Run()
}
