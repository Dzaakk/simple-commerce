package handler

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
	usecase "Dzaakk/simple-commerce/internal/transaction/usecase"
	template "Dzaakk/simple-commerce/package/template"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Usecase usecase.TransactionUseCase
}

func NewTransactionHandler(usecase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		Usecase: usecase,
	}
}

func (handler *TransactionHandler) Checkout(ctx *gin.Context) {
	var data model.TransactionReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return
	}

	fmt.Printf("input = %v", data)
	template.AuthorizedChecker(ctx, data.CustomerId)
	if ctx.IsAborted() {
		return
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
		return
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", receipt))
}
