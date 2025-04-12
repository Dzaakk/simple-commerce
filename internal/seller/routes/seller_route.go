package routes

import (
	handler "Dzaakk/simple-commerce/internal/seller/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type SellerRoutes struct {
	Handler *handler.SellerHandler
}

func NewSellerRoutes(handler *handler.SellerHandler) *SellerRoutes {
	return &SellerRoutes{Handler: handler}
}

func (sr *SellerRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	sellerHandler := r.Group("/api/v1/sellers")

	sellerHandler.Use()
	{
		sellerHandler.GET("", sr.Handler.FindAll)
		sellerHandler.GET("/id", sr.Handler.FindById)
		sellerHandler.GET("/username", sr.Handler.FindByUsername)
	}
}
