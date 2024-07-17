package main

import (
	psql "Dzaakk/synapsis/package/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db := psql.DB()
	r := gin.Default()
	fmt.Println("START")

	r.Run()
}
