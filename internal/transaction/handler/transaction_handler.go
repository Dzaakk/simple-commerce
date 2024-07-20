package transaction

import (
	model "Dzaakk/synapsis/internal/transaction/models"
	usecase "Dzaakk/synapsis/internal/transaction/usecase"
	auth "Dzaakk/synapsis/package/auth"
	template "Dzaakk/synapsis/package/template"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TransactionHandler struct {
	Usecase usecase.TransactionUseCase
}

func NewTransactionHandler(usecase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		Usecase: usecase,
	}
}

func (handler *TransactionHandler) Route(r *gin.RouterGroup, redis *redis.Client) {
	transactionHandler := r.Group("api/v1")

	transactionHandler.Use()
	{
		transactionHandler.POST("/transaction", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.Checkout(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "internal server error", err.Error()))
				return
			}
		})
	}
}

func (handler *TransactionHandler) Checkout(ctx *gin.Context) error {
	var data model.TransactionReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return nil
	}
	fmt.Printf("input = %v", data)

	template.AuthorizedChecker(ctx, data.CustomerId)
	if ctx.IsAborted() {
		return nil
	}

	receipt, err := handler.Usecase.CreateTransaction(data)
	if err != nil {
		var statusCode int
		var message string
		if err.Error() == "insufficient balance" {
			statusCode = http.StatusBadRequest
			message = "Insufficient Balance"
		} else {
			statusCode = http.StatusInternalServerError
			message = "Internal Server Error"
		}

		ctx.JSON(statusCode, template.Response(statusCode, message, err.Error()))
		return nil
	}
	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", receipt))
	return nil
}
