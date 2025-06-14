package route

import (
	"Dzaakk/simple-commerce/internal/transaction/handler"

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
		transactionHandler.POST("/transaction", tr.Handler.Checkout)
	}
}
