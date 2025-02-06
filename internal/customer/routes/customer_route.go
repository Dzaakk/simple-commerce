package routes

import (
	handler "Dzaakk/simple-commerce/internal/customer/handlers"
	"Dzaakk/simple-commerce/package/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CustomerRoutes struct {
	Handler *handler.CustomerHandler
}

func NewCustomerRoutes(handler *handler.CustomerHandler) *CustomerRoutes {
	return &CustomerRoutes{
		Handler: handler,
	}
}

func (cr *CustomerRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	customerHandler := r.Group("/api/v1")

	customerHandler.Use()
	{
		customerHandler.GET("/customers", auth.JWTMiddleware(redis), cr.Handler.FindCustomerById)
		customerHandler.POST("/balance", auth.JWTMiddleware(redis), cr.Handler.UpdateBalance)
	}
}
