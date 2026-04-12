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
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	err := h.service.RegisterCustomer(ctx, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success("Success Create Customer"))
}

func (h *AuthHandler) RegisterSeller(ctx *gin.Context) {
	var data dto.RegisterSellerRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	err := h.service.RegisterSeller(ctx, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success("Success Create Seller"))
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {

	activationCode := ctx.Query("code")
	if activationCode == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}
	err := h.service.VerifyEmail(ctx, activationCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Email verified successfully"))
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	req := &dto.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
		UserType: body.UserType,
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(res))
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var body dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	res, err := h.service.RefreshToken(ctx, body.RefreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(res))
}
func (h *AuthHandler) Logout(ctx *gin.Context) {
	var body dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.service.Logout(ctx, body.RefreshToken); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Logged out successfully"))
}
