package handler

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	sellerUsecase "Dzaakk/simple-commerce/internal/seller/usecase"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthHandler struct {
	Usecase       usecase.AuthUseCase
	SellerUsecase sellerUsecase.SellerUseCase
}

func NewAtuhHandler(usecase usecase.AuthUseCase, sellerUsecase sellerUsecase.SellerUseCase) *AuthHandler {
	return &AuthHandler{
		Usecase:       usecase,
		SellerUsecase: sellerUsecase,
	}
}

func (h *AuthHandler) RegistrationCustomer(ctx *gin.Context) {
	var data model.CustomerRegistrationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.Usecase.RegistrationCustomer(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Create User"))
}

func (h *AuthHandler) ActivationCustomer(ctx *gin.Context) {
	var data model.CustomerActivationReq

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

func (h *AuthHandler) RegistrationSeller(ctx *gin.Context) {
	var data model.SellerRegistrationReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	_, err := h.Usecase.SellerRegistration(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success("Success Create Seller"))
}
func (h *AuthHandler) ActivationSeller(ctx *gin.Context, r *redis.Client) {

}

func (h *AuthHandler) LoginSeller(ctx *gin.Context) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	data, err := h.SellerUsecase.FindByEmail(ctx, reqData.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}

	if !util.CheckPasswordHash(reqData.Password, data.Password) {
		ctx.JSON(http.StatusBadRequest, response.InvalidEmailOrPassword())
		return
	}
	// byteData, err := json.Marshal(data)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
	// 	return
	// }
	// _ = auth.NewTokenGenerator(r, byteData)

	ctx.JSON(http.StatusOK, response.Success("Login Success"))
}
