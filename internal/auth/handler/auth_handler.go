package handler

import (
	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/service"
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(usecase service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: usecase,
	}
}

func (h *AuthHandler) RegisterCustomer(ctx *gin.Context) {
	var data dto.RegisterCustomerRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.service.RegisterCustomer(ctx, &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success("Success Create Customer"))
}

func (h *AuthHandler) RegisterSeller(ctx *gin.Context) {
	var data dto.RegisterSellerRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.service.RegisterSeller(ctx, &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success("Success Create Seller"))
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {

	activationCode := ctx.Query("code")
	if activationCode == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}
	err := h.service.VerifyEmail(ctx, activationCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Email verified successfully"))
}
