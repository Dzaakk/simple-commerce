package handler

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	emailUsecase "Dzaakk/simple-commerce/internal/email/usecase"
	sellerUsecase "Dzaakk/simple-commerce/internal/seller/usecase"
	"Dzaakk/simple-commerce/package/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase       usecase.AuthUseCase
	SellerUsecase sellerUsecase.SellerUseCase
	EmailUsecase  emailUsecase.EmailUseCase
}

func NewAtuhHandler(usecase usecase.AuthUseCase, sellerUsecase sellerUsecase.SellerUseCase, emailUsecase emailUsecase.EmailUseCase) *AuthHandler {
	return &AuthHandler{
		Usecase:       usecase,
		SellerUsecase: sellerUsecase,
		EmailUsecase:  emailUsecase,
	}
}

func (h *AuthHandler) RegistrationCustomer(ctx *gin.Context) {
	var data model.CustomerRegistrationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	dataEmail, err := h.Usecase.RegistrationCustomer(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	err = h.EmailUsecase.SendEmailActivation(ctx, *dataEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Create User"))
}

func (h *AuthHandler) ActivationCustomer(ctx *gin.Context) {
	var data model.ActivationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.Usecase.ActivationCustomer(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Activate User"))
}

func (h *AuthHandler) LoginCustomer(ctx *gin.Context) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	err := h.Usecase.LoginCustomer(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Login Success"))
}
func (h *AuthHandler) ForgotPasswordCustomer(ctx *gin.Context) {}

func (h *AuthHandler) RegistrationSeller(ctx *gin.Context) {
	var data model.SellerRegistrationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	dataEmail, err := h.Usecase.RegistrationSeller(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	err = h.EmailUsecase.SendEmailActivation(ctx, *dataEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Create Seller"))
}

func (h *AuthHandler) ActivationSeller(ctx *gin.Context) {

	var data model.ActivationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.Usecase.ActivationSeller(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Activate User"))
}

func (h *AuthHandler) LoginSeller(ctx *gin.Context) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	err := h.Usecase.LoginSeller(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Login Success"))
}
func (h *AuthHandler) ForgotPasswordSeller(ctx *gin.Context) {}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	roleRaw, roleExists := ctx.Get("role")
	emailRaw, emailExists := ctx.Get("email")

	role, okRole := roleRaw.(string)
	email, okEmail := emailRaw.(string)

	if !roleExists || !emailExists || !okRole || !okEmail {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(""))
		return
	}

	err := h.Usecase.Logout(ctx, role, email)
	if err != nil {
		log.Printf("[Logout] Failed to delete token for %s (%s): %v", role, email, err)
	}

	ctx.JSON(http.StatusOK, response.Success("Logged out successfully"))
}
