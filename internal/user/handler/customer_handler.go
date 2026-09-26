package handler

import (
	"net/http"

	"Dzaakk/simple-commerce/internal/middleware"
	"Dzaakk/simple-commerce/internal/user/service"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	customers service.CustomerService
}

func NewUserHandler(customers service.CustomerService) *UserHandler {
	return &UserHandler{customers: customers}
}

func (h *UserHandler) Me(ctx *gin.Context) {
	value, exists := ctx.Get(middleware.UserIDKey)
	userID, ok := value.(string)
	if !exists || !ok || userID == "" {
		ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
		return
	}
	customer, err := h.customers.FindByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}
	if customer == nil {
		ctx.Error(response.ErrUserNotFound)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(customer))
}
