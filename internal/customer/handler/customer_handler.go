package customer

import (
	model "Dzaakk/synapsis/internal/customer/models"
	usecase "Dzaakk/synapsis/internal/customer/usecase"
	utils "Dzaakk/synapsis/package/template"
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

func (handler *CustomerHandler) Route(r *gin.RouterGroup) {
	customerHandler := r.Group("/api/v1")

	customerHandler.Use()
	{
		customerHandler.POST("/customers", func(ctx *gin.Context) {
			if err := handler.Create(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "failed to create customer", err.Error()))
				return
			}
		})
		customerHandler.GET("/customers", func(ctx *gin.Context) {
			if err := handler.FindCustomerById(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "failed to create customer", err.Error()))
				return
			}
		})
		customerHandler.POST("/balance", func(ctx *gin.Context) {
			if err := handler.UpdateBalance(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "failed to create customer", err.Error()))
				return
			}
		})
		// customerHandler.GET("/balance", func(ctx *gin.Context) {
		// 	if err := handler.GetBalance(ctx); err != nil {
		// 		ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "failed to create customer", err.Error()))
		// 		return
		// 	}
		// })
	}
}

func (handler *CustomerHandler) FindCustomerById(ctx *gin.Context) error {
	id, errParam := strconv.Atoi(ctx.Query("id"))
	if errParam != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, "missing required parameter 'id'", errParam.Error()))
		return nil
	}

	data, err := handler.Usecase.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.Response(http.StatusNotFound, "not found", err.Error()))
		ctx.Abort()
		return nil
	}

	ctx.JSON(http.StatusOK, utils.Response(http.StatusOK, "Success", data))
	return nil
}
func (handler *CustomerHandler) Create(ctx *gin.Context) error {
	var data model.CustomerReq

	if err := ctx.ShouldBindJSON(&data); err != nil {
		newErr := fmt.Errorf("failed to bind JSON: %s", err)

		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return newErr
	}
	fmt.Printf("input = %v", data)
	id, err := handler.Usecase.Create(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return err
	}
	res := model.CustomerRes{
		Id: fmt.Sprintf("%d", *id),
	}
	ctx.JSON(http.StatusCreated, utils.Response(http.StatusCreated, "Success Create User", res))
	return nil
}

func (handler *CustomerHandler) UpdateBalance(ctx *gin.Context) error {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, "missing required parameter 'id'", err.Error()))
		return nil
	}
	balance, err := strconv.Atoi(ctx.Query("balance"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, "missing required parameter 'balance'", err.Error()))
		return nil
	}
	actionType := ctx.Query("actionType")
	oldData, err := handler.Usecase.GetBalance(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return err
	}

	newBalance, err := handler.Usecase.UpdateBalance(id, float32(balance), actionType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.Response(http.StatusNotFound, "not found", err.Error()))
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
	ctx.JSON(http.StatusOK, utils.Response(http.StatusOK, "Success Update Balance", response))
	return nil
}
