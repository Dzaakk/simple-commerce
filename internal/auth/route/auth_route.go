package route

import (
	"Dzaakk/simple-commerce/internal/auth/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	Handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		Handler: handler,
	}
}

func (ar *AuthRoutes) Route(r *gin.RouterGroup) {
	apiGroup := r.Group("/api/v1")

	apiGroup.Use()
	{
		apiGroup.POST("/customer-register", ar.Handler.RegistrationCustomer)
		apiGroup.POST("/customer-activation", func(ctx *gin.Context) {
			ar.Handler.ActivationCustomer(ctx)
		})
		apiGroup.POST("/customer-login", func(ctx *gin.Context) {
			ar.Handler.LoginCustomer(ctx)
		})

		apiGroup.POST("/register-seller", ar.Handler.RegistrationSeller)
		apiGroup.POST("/login-seller", func(ctx *gin.Context) {
			ar.Handler.LoginSeller(ctx)
		})
	}
}
