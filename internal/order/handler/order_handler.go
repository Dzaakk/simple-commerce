package handler

import (
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/service"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if req.CustomerID == "" {
		if id, ok := ctx.Get("id"); ok {
			if idStr, ok := id.(string); ok {
				req.CustomerID = idStr
			}
		}
	}

	res, err := h.Service.CreateOrder(ctx, &req)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(res))
}

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	filter, err := parseOrderFilter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.Service.GetOrdersByCustomer(ctx, customerID, filter)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *OrderHandler) GetOrderDetail(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.Service.GetOrderByID(ctx, customerID, orderID)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}
	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NotFound("order not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	customerID, ok := getCustomerID(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.Service.CancelOrder(ctx, customerID, orderID); err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Cancel Order"))
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

func parseOrderFilter(ctx *gin.Context) (dto.OrderFilter, error) {
	var filter dto.OrderFilter

	if status := ctx.Query("status"); status != "" {
		st := constant.OrderStatus(status)
		filter.Status = &st
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return filter, err
		}
		filter.Page = page
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return filter, err
		}
		filter.Limit = limit
	}

	return filter, nil
}

func writeOrderError(ctx *gin.Context, err error) {
	msg := err.Error()
	if strings.Contains(msg, "invalid") {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(msg))
		return
	}
	if strings.Contains(msg, "unauthorized") {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(msg))
		return
	}
	if strings.Contains(msg, "not found") {
		ctx.JSON(http.StatusNotFound, response.NotFound(msg))
		return
	}
	if strings.Contains(msg, "stock") {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(msg))
		return
	}

	ctx.JSON(http.StatusInternalServerError, response.InternalServerError(msg))
}
