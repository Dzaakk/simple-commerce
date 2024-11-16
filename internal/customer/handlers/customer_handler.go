package handler

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	usecase "Dzaakk/simple-commerce/internal/customer/usecases"
	auth "Dzaakk/simple-commerce/package/auth"
	db "Dzaakk/simple-commerce/package/db"
	template "Dzaakk/simple-commerce/package/templates"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CustomerHandler struct {
	Usecase usecase.CustomerUseCase
}

func NewCustomerHandler(usecase usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		Usecase: usecase,
	}
}

func (handler *CustomerHandler) Login(ctx *gin.Context, redis *redis.Client) {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid email or password"))
		return
	}

	data, err := handler.Usecase.FindByEmail(reqData.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	if !template.CheckPasswordHash(reqData.Password, data.Password) {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid email or password"))
		return
	}

	_, err = auth.TokenJWTGenerator(db.Redis(), *data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusInternalServerError, "Internal Server Error", "Error generate JWT token"))
		return
	}
	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Login Success", "You have 30 mins token Session!"))
	return
}

func (handler *CustomerHandler) FindCustomerById(ctx *gin.Context) {
	id, errParam := strconv.Atoi(ctx.Query("id"))
	if errParam != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "missing required parameter 'id'", errParam.Error()))
		return
	}

	template.AuthorizedChecker(ctx, ctx.Query("id"))

	data, err := handler.Usecase.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", err.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success", data))
	return
}

func (handler *CustomerHandler) Create(ctx *gin.Context) {
	var data model.CustomerReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return
	}
	fmt.Printf("input = %v", data)
	id, err := handler.Usecase.Create(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	res := model.CustomerRes{
		Id: fmt.Sprintf("%d", *id),
	}
	ctx.JSON(http.StatusCreated, template.Response(http.StatusCreated, "Success Create User", res))
	return
}

func (handler *CustomerHandler) UpdateBalance(ctx *gin.Context) {
	var data model.BalanceUpdateReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return
	}

	balance, err := strconv.ParseFloat(data.Balance, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data on field 'balance'"))
		return
	}

	id, _ := strconv.Atoi(data.Id)
	oldData, err := handler.Usecase.GetBalance(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	template.AuthorizedChecker(ctx, data.Id)
	newBalance, err := handler.Usecase.UpdateBalance(id, float64(balance), data.ActionType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", err.Error()))
		ctx.Abort()
		return
	}

	response := model.BalanceUpdateRes{
		BalanceOld: *oldData,
		BalanceNew: model.CustomerBalance{
			Id:      id,
			Balance: *newBalance,
		},
	}
	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success Update Balance", response))
	return
}
