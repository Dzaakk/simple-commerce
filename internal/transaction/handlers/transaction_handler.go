package handler

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
	usecase "Dzaakk/simple-commerce/internal/transaction/usecases"
	"Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/templates"
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
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	template.AuthorizedChecker(ctx, data.CustomerID)
	if ctx.IsAborted() {
		return
	}

	receipt, err := handler.Usecase.CreateTransaction(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(receipt))
}
