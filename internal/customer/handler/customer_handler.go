package handler

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
	"fmt"
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

func (h *CustomerHandler) FindCustomerByUsername(ctx *gin.Context) {}
func (h *CustomerHandler) Update(ctx *gin.Context)                 {}

func (h *CustomerHandler) ChangePassword(ctx *gin.Context) {
	var req model.ChangePasswordReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	customerID, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	template.AuthorizedChecker(ctx, ctx.Query("id"))

	_, err = h.Usecase.UpdatePassword(ctx, customerID, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success change password!"))
}

func (h *CustomerHandler) UpdateBalance(ctx *gin.Context) {
	var data model.BalanceUpdateReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	balance, err := strconv.ParseFloat(data.Balance, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	id, _ := strconv.ParseInt(data.CustomerID, 10, 64)
	oldData, err := h.Usecase.GetBalance(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	template.AuthorizedChecker(ctx, data.CustomerID)
	newBalance, err := h.Usecase.UpdateBalance(ctx, id, float64(balance), data.ActionType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.NotFound(err.Error()))
		ctx.Abort()
		return
	}

	res := model.BalanceUpdateRes{
		BalanceOld: *oldData,
		BalanceNew: model.CustomerBalanceRes{
			CustomerID: data.CustomerID,
			Balance:    fmt.Sprintf("%.2f", newBalance),
		},
	}
	ctx.JSON(http.StatusOK, response.Success(res))
}
