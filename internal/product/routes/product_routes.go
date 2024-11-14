package routes

import (
	handler "Dzaakk/simple-commerce/internal/product/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ProductRoutes struct {
	Handler *handler.ProductHandler
}

func NewProductRoutes(handler *handler.ProductHandler) *ProductRoutes {
	return &ProductRoutes{
		Handler: handler,
	}
}
func (pr *ProductRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	productHandler := r.Group("api/v1")

	productHandler.Use()
	{
		productHandler.GET("/product", pr.Handler.GetProduct)
	}
}
