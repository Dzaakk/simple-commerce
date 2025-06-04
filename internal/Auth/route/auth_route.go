package route

import (
	"Dzaakk/simple-commerce/internal/auth/handler"

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
	apiGroup := r.Group("/api/v1")

	apiGroup.Use()
	{
		apiGroup.POST("/register-customer", ar.Handler.RegistrationCustomer)
		apiGroup.POST("/login-customer", func(ctx *gin.Context) {
			ar.Handler.LoginCustomer(ctx, redis)
		})

		apiGroup.POST("/register-seller", ar.Handler.RegistrationSeller)
		apiGroup.POST("/login-seller", func(ctx *gin.Context) {
			ar.Handler.LoginSeller(ctx, redis)
		})
	}
}
