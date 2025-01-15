package handlers

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	usecase "Dzaakk/simple-commerce/internal/customer/usecases"
	"Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/templates"
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

func (handler *CustomerHandler) Login(ctx *gin.Context) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid email or password"))
		return
	}

	data, err := handler.Usecase.FindByEmail(ctx, reqData.Email)
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success("Login Success"))
}

func (handler *CustomerHandler) FindCustomerById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	template.AuthorizedChecker(ctx, ctx.Query("id"))

	data, err := handler.Usecase.FindById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NotFound(err.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (handler *CustomerHandler) Create(ctx *gin.Context) {
	var data model.CreateReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}
	_, err := handler.Usecase.Create(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success("Success Create User"))
}

func (handler *CustomerHandler) UpdateBalance(ctx *gin.Context) {
	var data model.BalanceUpdateReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}

	balance, err := strconv.ParseFloat(data.Balance, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	id, _ := strconv.ParseInt(data.Id, 10, 64)
	oldData, err := handler.Usecase.GetBalance(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	template.AuthorizedChecker(ctx, data.Id)
	newBalance, err := handler.Usecase.UpdateBalance(ctx, id, float64(balance), data.ActionType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.NotFound(err.Error()))
		ctx.Abort()
		return
	}

	res := model.BalanceUpdateRes{
		BalanceOld: *oldData,
		BalanceNew: model.CustomerBalanceRes{
			Id:      data.Id,
			Balance: fmt.Sprintf("%.2f", newBalance),
		},
	}
	ctx.JSON(http.StatusOK, response.Success(res))
}
