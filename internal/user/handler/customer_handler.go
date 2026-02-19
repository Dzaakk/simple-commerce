package handler

import (
	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/service"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.CustomerService
}

func NewUserHandler(service service.CustomerService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) FindCustomerByID(ctx *gin.Context) {

	customerID := ctx.Query("id")
	if !util.AuthorizedChecker(ctx, ctx.Query("id")) {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(""))
		return
	}

	data, err := h.Service.FindByID(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NotFound("user not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *UserHandler) Update(ctx *gin.Context) {
	var req dto.UpdateReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update!"))
}
