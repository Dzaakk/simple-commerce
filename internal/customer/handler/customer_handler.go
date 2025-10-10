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
	Usecase usecase.CustomerUseCase
}

func NewCustomerHandler(usecase usecase.CustomerUseCase) *CustomerHandler {
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

	template.AuthorizedChecker(ctx, ctx.Query("id"))

	data, err := h.Usecase.FindByID(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NotFound(err.Error()))
		ctx.Abort()
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

	err := h.Usecase.Update(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update!"))
}
