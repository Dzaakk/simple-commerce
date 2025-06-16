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

	customer := apiGroup.Group("/customer")
	{
		customer.POST("/register", ar.Handler.RegistrationCustomer)
		customer.POST("/activate", ar.Handler.ActivationCustomer)
		customer.POST("/login", ar.Handler.LoginCustomer)
	}

	seller := apiGroup.Group("/seller")
	{
		seller.POST("/register", ar.Handler.RegistrationSeller)
		seller.POST("/activate", ar.Handler.ActivationSeller)
		seller.POST("/login", ar.Handler.LoginSeller)
	}
}
