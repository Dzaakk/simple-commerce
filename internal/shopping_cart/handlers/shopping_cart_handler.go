package handler

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecases"
	template "Dzaakk/simple-commerce/package/templates"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShoppingCartHandler struct {
	Usecase usecase.ShoppingCartUseCase
}

func NewShoppingCartHandler(usecase usecase.ShoppingCartUseCase) *ShoppingCartHandler {
	return &ShoppingCartHandler{
		Usecase: usecase,
	}
}

func (handler *ShoppingCartHandler) AddProductToShoppingCart(ctx *gin.Context) {
	var reqData model.ShoppingCartReq
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return
	}

	template.AuthorizedChecker(ctx, reqData.CustomerId)
	newShopingCart, err := handler.Usecase.Add(reqData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", newShopingCart))
}

func (handler *ShoppingCartHandler) GetListShoppingCart(ctx *gin.Context) {
	customerId, _ := strconv.Atoi(ctx.Query("customerId"))

	template.AuthorizedChecker(ctx, ctx.Query("customerId"))
	if ctx.IsAborted() {
		return
	}
	listShoppingCart, err := handler.Usecase.GetListItem(customerId)
	if err != nil {

		var statusCode int
		var message string
		if err.Error() == "cart is empty" {
			statusCode = http.StatusOK
			message = "your shopping cart is empty."
		} else {
			statusCode = http.StatusInternalServerError
			message = "Internal Server Error"
		}

		ctx.JSON(statusCode, template.Response(statusCode, message, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", listShoppingCart))
}

func (handler *ShoppingCartHandler) DeleteShoppingList(ctx *gin.Context) {
	var data model.DeleteReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return
	}
	fmt.Printf("input = %v", data)
	template.AuthorizedChecker(ctx, data.CustomerId)
	if ctx.IsAborted() {
		return
	}

	err := handler.Usecase.DeleteShoppingList(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", "Success Delete Shopping List"))
}
