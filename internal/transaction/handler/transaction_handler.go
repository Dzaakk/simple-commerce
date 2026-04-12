package handler

import (
	"Dzaakk/simple-commerce/internal/transaction/dto"
	"Dzaakk/simple-commerce/internal/transaction/service"
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	var req dto.CreateTransactionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if req.CustomerID == "" {
		if id, ok := ctx.Get("id"); ok {
			if idStr, ok := id.(string); ok {
				req.CustomerID = idStr
			}
		}
	}

	res, err := h.Service.CreateTransaction(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(res))
}

func (h *TransactionHandler) GetTransactionByID(ctx *gin.Context) {
	transactionID := ctx.Param("id")
	if transactionID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	customerID, ok := getCustomerID(ctx)
	if !ok {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	res, err := h.Service.GetTransactionByID(ctx, customerID, transactionID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(res))
}

func (h *TransactionHandler) GetTransactionByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	if orderID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	customerID, ok := getCustomerID(ctx)
	if !ok {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	res, err := h.Service.GetTransactionByOrderID(ctx, customerID, orderID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(res))
}

func (h *TransactionHandler) PaymentCallback(ctx *gin.Context) {
	var req dto.PaymentCallbackReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.Service.HandlePaymentCallback(ctx, &req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("callback processed"))
}

func (h *TransactionHandler) ExpireTransaction(ctx *gin.Context) {
	transactionID := ctx.Param("id")
	if transactionID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.Service.ExpireTransaction(ctx, transactionID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("transaction expired"))
}

func getCustomerID(ctx *gin.Context) (string, bool) {
	if id, ok := ctx.Get("id"); ok {
		if idStr, ok := id.(string); ok && idStr != "" {
			if q := ctx.Query("customer_id"); q != "" && q != idStr {
				return "", false
			}
			return idStr, true
		}
	}

	if q := ctx.Query("customer_id"); q != "" {
		return q, true
	}

	return "", false
}
