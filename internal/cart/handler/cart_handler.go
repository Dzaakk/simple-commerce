package handler

import (
	"Dzaakk/simple-commerce/internal/cart/service"
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	Service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{Service: service}
}

type cartItemReq struct {
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int    `json:"quantity"`
}

func (h *CartHandler) GetCart(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx, "")
	if !ok {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.Service.GetCartItems(ctx, customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) AddItem(ctx *gin.Context) {
	var req cartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	customerID, ok := getCustomerID(ctx, req.CustomerID)
	if !ok {
		ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
		return
	}

	data, err := h.Service.AddItem(ctx, customerID, req.ProductID, req.Quantity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) UpdateItem(ctx *gin.Context) {
	var req cartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	customerID, ok := getCustomerID(ctx, req.CustomerID)
	if !ok {
		ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
		return
	}

	data, err := h.Service.UpdateItem(ctx, customerID, req.ProductID, req.Quantity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) DeleteItem(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	customerID, ok := getCustomerID(ctx, "")
	if !ok {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.Service.DeleteItem(ctx, customerID, productID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Delete Item"))
}

func (h *CartHandler) ClearItems(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx, "")
	if !ok {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.Service.ClearItems(ctx, customerID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Clear Items"))
}

func getCustomerID(ctx *gin.Context, bodyCustomerID string) (string, bool) {
	if idVal, exists := ctx.Get("id"); exists {
		if id, ok := idVal.(string); ok && id != "" {
			if bodyCustomerID != "" && bodyCustomerID != id {
				return "", false
			}
			if queryID := ctx.Query("customer_id"); queryID != "" && queryID != id {
				return "", false
			}
			return id, true
		}
	}

	if bodyCustomerID != "" {
		return bodyCustomerID, true
	}
	if queryID := ctx.Query("customer_id"); queryID != "" {
		return queryID, true
	}

	return "", false
}
