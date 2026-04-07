package handler

import (
	"Dzaakk/simple-commerce/internal/cart/service"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strings"

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
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.Service.GetCartItems(ctx, customerID)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) AddItem(ctx *gin.Context) {
	var req cartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	customerID, ok := getCustomerID(ctx, req.CustomerID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(""))
		return
	}

	data, err := h.Service.AddItem(ctx, customerID, req.ProductID, req.Quantity)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) UpdateItem(ctx *gin.Context) {
	var req cartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	customerID, ok := getCustomerID(ctx, req.CustomerID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(""))
		return
	}

	data, err := h.Service.UpdateItem(ctx, customerID, req.ProductID, req.Quantity)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CartHandler) DeleteItem(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	customerID, ok := getCustomerID(ctx, "")
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.Service.DeleteItem(ctx, customerID, productID); err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Delete Item"))
}

func (h *CartHandler) ClearItems(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx, "")
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.Service.ClearItems(ctx, customerID); err != nil {
		writeError(ctx, err)
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

func writeError(ctx *gin.Context, err error) {
	msg := err.Error()
	if strings.HasPrefix(msg, "invalid parameter") || strings.Contains(msg, "stock") {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(msg))
		return
	}
	if strings.Contains(msg, "not found") {
		ctx.JSON(http.StatusNotFound, response.NotFound(msg))
		return
	}

	ctx.JSON(http.StatusInternalServerError, response.InternalServerError(msg))
}
