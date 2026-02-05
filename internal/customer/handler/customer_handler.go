package handler

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	Usecase usecase.CustomerUsecase
}

func NewCustomerHandler(usecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		Usecase: usecase,
	}
}

func (h *CustomerHandler) FindCustomerByID(ctx *gin.Context) {
	customerID, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if !template.AuthorizedChecker(ctx, ctx.Query("id")) {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(""))
		return
	}

	data, err := h.Usecase.FindByID(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NotFound("customer not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CustomerHandler) Update(ctx *gin.Context) {
	var req model.UpdateReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	err := h.Usecase.Update(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update!"))
}
