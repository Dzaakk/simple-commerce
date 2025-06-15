package route

import (
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"Dzaakk/simple-commerce/internal/transaction/handler"

	"github.com/gin-gonic/gin"
)

type TransactionRoutes struct {
	Handler       *handler.TransactionHandler
	JWTMiddleware *middleware.JWTMiddleware
}

func NewTransactionRoutes(handler *handler.TransactionHandler, jwtMiddleware *middleware.JWTMiddleware) *TransactionRoutes {
	return &TransactionRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}

func (tr *TransactionRoutes) Route(r *gin.RouterGroup) {
	transactionHandler := r.Group("api/v1")

	transactionHandler.Use(tr.JWTMiddleware.ValidateToken())
	{
		transactionHandler.POST("/transaction", tr.Handler.Checkout)
	}
}
