package routes

import (
	handler "Dzaakk/simple-commerce/internal/auth/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthRoutes struct {
	Handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		Handler: handler,
	}
}

func (ar *AuthRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	authHandler := r.Group("/api/v1")

	authHandler.Use()
	{
		authHandler.POST("/login-customer", ar.Handler.LoginCustomer)
		authHandler.POST("/register-customer", ar.Handler.RegistrationCustomer)
		authHandler.POST("/login-seller", ar.Handler.LoginSeller)
		authHandler.POST("/register-seller", ar.Handler.RegistrationSeller)
	}
}
