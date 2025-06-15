package route

import (
	"Dzaakk/simple-commerce/internal/customer/handler"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type CustomerRoutes struct {
	Handler       *handler.CustomerHandler
	JWTMiddleware *middleware.JWTMiddleware
}

func NewCustomerRoutes(handler *handler.CustomerHandler, jwtMiddleware *middleware.JWTMiddleware) *CustomerRoutes {
	return &CustomerRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}

func (cr *CustomerRoutes) Route(r *gin.RouterGroup) {
	customerHandler := r.Group("/api/v1")
	customerHandler.Use(cr.JWTMiddleware.ValidateToken())
	{
		customerHandler.GET("/customers", cr.Handler.FindCustomerById)
		customerHandler.POST("/balance", cr.Handler.UpdateBalance)
	}
}
