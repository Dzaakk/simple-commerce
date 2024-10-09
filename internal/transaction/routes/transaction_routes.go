package transaction

import (
	handler "Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/package/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Route(r *gin.RouterGroup, redis *redis.Client) {
	transactionHandler := r.Group("api/v1")

	transactionHandler.Use()
	{
		transactionHandler.POST("/transaction", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			handler.Checkout(ctx)
		})
	}
}
