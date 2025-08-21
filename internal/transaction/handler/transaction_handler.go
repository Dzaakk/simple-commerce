package handler

import (
	"Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/internal/transaction/usecase"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
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

func (h *TransactionHandler) Checkout(ctx *gin.Context) {
	var data model.TransactionReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	template.AuthorizedChecker(ctx, data.CustomerID)
	if ctx.IsAborted() {
		return
	}

	receipt, err := h.Usecase.CreateTransaction(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(receipt))
}

func (h *TransactionHandler) GetListHistoryTransaction(ctx *gin.Context) {}
func (h *TransactionHandler) Paid(ctx *gin.Context)                      {}
