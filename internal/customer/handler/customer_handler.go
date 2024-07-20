package customer

import (
	model "Dzaakk/synapsis/internal/customer/models"
	usecase "Dzaakk/synapsis/internal/customer/usecase"
	auth "Dzaakk/synapsis/package/auth"
	db "Dzaakk/synapsis/package/db"
	template "Dzaakk/synapsis/package/template"
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

func (handler *CustomerHandler) Route(r *gin.RouterGroup, redis *redis.Client) {
	customerHandler := r.Group("/api/v1")

	customerHandler.Use()
	{
		customerHandler.POST("/login", func(ctx *gin.Context) {
			if err := handler.Login(ctx, redis); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "login failed", err.Error()))
				return
			}
		})
		customerHandler.POST("/register", func(ctx *gin.Context) {
			if err := handler.Create(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to create customer", err.Error()))
				return
			}
		})
		customerHandler.GET("/customers", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.FindCustomerById(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to Get customer", err.Error()))
				return
			}
		})
		customerHandler.POST("/balance", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.UpdateBalance(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to update balance", err.Error()))
				return
			}
		})
	}
}

func (handler *CustomerHandler) Login(ctx *gin.Context, redis *redis.Client) error {
	var reqData model.LoginReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid email or password"))
		return nil
	}

	data, err := handler.Usecase.FindByEmail(reqData.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return nil
	}
	if !template.CheckPasswordHash(reqData.Password, data.Password) {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid email or password"))
		return nil
	}

	_, err = auth.TokenJWTGenerator(db.Redis(), *data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusInternalServerError, "Internal Server Error", "Error generate JWT token"))
		return nil
	}
	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Login Success", "You have 30 mins token Session!"))
	return nil
}

func (handler *CustomerHandler) FindCustomerById(ctx *gin.Context) error {
	id, errParam := strconv.Atoi(ctx.Query("id"))
	if errParam != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "missing required parameter 'id'", errParam.Error()))
		return nil
	}

	template.AuthorizedChecker(ctx, ctx.Query("id"))

	data, err := handler.Usecase.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", err.Error()))
		ctx.Abort()
		return nil
	}

	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success", data))
	return nil
}

func (handler *CustomerHandler) Create(ctx *gin.Context) error {
	var data model.CustomerReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return nil
	}
	fmt.Printf("input = %v", data)
	id, err := handler.Usecase.Create(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return nil
	}
	res := model.CustomerRes{
		Id: fmt.Sprintf("%d", *id),
	}
	ctx.JSON(http.StatusCreated, template.Response(http.StatusCreated, "Success Create User", res))
	return nil
}

func (handler *CustomerHandler) UpdateBalance(ctx *gin.Context) error {
	var data model.BalanceUpdateReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return nil
	}

	balance, err := strconv.ParseFloat(data.Balance, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data on field 'balance'"))
		return nil
	}

	id, _ := strconv.Atoi(data.Id)
	oldData, err := handler.Usecase.GetBalance(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return nil
	}
	template.AuthorizedChecker(ctx, data.Id)
	newBalance, err := handler.Usecase.UpdateBalance(id, float32(balance), data.ActionType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", err.Error()))
		ctx.Abort()
		return nil
	}

	response := model.BalanceUpdateRes{
		BalanceOld: *oldData,
		BalanceNew: model.CustomerBalance{
			Id:      id,
			Balance: *newBalance,
		},
	}
	ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success Update Balance", response))
	return nil
}
