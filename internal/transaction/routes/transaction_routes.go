package routes

import (
	handler "Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/package/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TransactionRoutes struct {
	Handler *handler.TransactionHandler
}

func NewTransactionRoutes(handler *handler.TransactionHandler) *TransactionRoutes {
	return &TransactionRoutes{
		Handler: handler,
	}
}

func (tr *TransactionRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	transactionHandler := r.Group("api/v1")

	transactionHandler.Use()
	{
		transactionHandler.POST("/transaction", auth.JWTMiddleware(redis), tr.Handler.Checkout)
	}
}
