package handler

import (
	"net/http"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/service"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) RegisterCustomer(ctx *gin.Context) {
	var req dto.RegisterCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}
	if err := h.service.RegisterCustomer(ctx.Request.Context(), &req); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, response.Response(http.StatusCreated, "Created", nil))
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}
	result, err := h.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(result))
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}
	result, err := h.service.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(result))
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}
	if err := h.service.Logout(ctx.Request.Context(), req.RefreshToken); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
