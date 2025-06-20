package route

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	Handler            *handler.AuthHandler
	CustomerMiddleware *middleware.JWTCustomerMiddleware
	SellerMiddleware   *middleware.JWTSellerMiddleware
}

func NewAuthRoutes(handler *handler.AuthHandler, customerMiddleware *middleware.JWTCustomerMiddleware, sellerMiddleware *middleware.JWTSellerMiddleware) *AuthRoutes {
	return &AuthRoutes{
		Handler:            handler,
		CustomerMiddleware: customerMiddleware,
		SellerMiddleware:   sellerMiddleware,
	}
}

func (ar *AuthRoutes) Route(r *gin.RouterGroup) {
	apiGroup := r.Group("/api/v1")

	customer := apiGroup.Group("/customer")
	{
		customer.POST("/register", ar.Handler.RegistrationCustomer)
		customer.POST("/activate", ar.Handler.ActivationCustomer)
		customer.POST("/login", ar.Handler.LoginCustomer)
		customer.POST("/logout", ar.CustomerMiddleware.ValidateToken(), ar.Handler.Logout)
	}

	seller := apiGroup.Group("/seller")
	{
		seller.POST("/register", ar.Handler.RegistrationSeller)
		seller.POST("/activate", ar.Handler.ActivationSeller)
		seller.POST("/login", ar.Handler.LoginSeller)
		// customer.POST("/logout", ar.SellerMiddleware.ValidateToken(), ar.Handler.Logout)
	}
}
