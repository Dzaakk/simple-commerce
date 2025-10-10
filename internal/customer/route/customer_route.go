package route

import (
	"Dzaakk/simple-commerce/internal/customer/handler"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type CustomerRoutes struct {
	Handler       *handler.CustomerHandler
	JWTMiddleware *middleware.JWTCustomerMiddleware
}

func NewCustomerRoutes(handler *handler.CustomerHandler, jwtMiddleware *middleware.JWTCustomerMiddleware) *CustomerRoutes {
	return &CustomerRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}

func (cr *CustomerRoutes) Route(r *gin.RouterGroup) {
	customerHandler := r.Group("/api/v1/customer")
	customerHandler.Use(cr.JWTMiddleware.ValidateToken())
	{
		customerHandler.GET("/find-all", cr.Handler.FindCustomerByID)
	}
}
