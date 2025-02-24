package handlers

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	usecase "Dzaakk/simple-commerce/internal/auth/usecases"
	custUsecase "Dzaakk/simple-commerce/internal/customer/usecases"
	"Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase         usecase.AuthUseCase
	CustomerUsecase custUsecase.CustomerUseCase
}

func NewAtuhHandler(usecase usecase.AuthUseCase, custUsecase custUsecase.CustomerUseCase) *AuthHandler {
	return &AuthHandler{
		Usecase:         usecase,
		CustomerUsecase: custUsecase,
	}
}

func (h *AuthHandler) RegistrationCustomer(ctx *gin.Context) {
	var data model.CustomerRegistration

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}
	_, err := h.Usecase.CustomerRegistration(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success("Success Create User"))
}

func (h *AuthHandler) LoginCustomer(ctx *gin.Context) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid email or password"))
		return
	}

	data, err := h.CustomerUsecase.FindByEmail(ctx, reqData.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	if !template.CheckPasswordHash(reqData.Password, data.Password) {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid email or password"))
		return
	}
	// cache token
	// _, err = auth.NewTokenGenerator(db.Redis(), *data)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
	// 	return
	// }
	ctx.JSON(http.StatusOK, response.Success("Login Success"))
}

func (h *AuthHandler) RegistrationSeller(ctx *gin.Context) {
	var data model.SellerRegistration

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}
	_, err := h.Usecase.SellerRegistration(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success("Success Create Seller"))
}
func (h *AuthHandler) LoginSeller(ctx *gin.Context) {
}
