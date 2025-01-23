package handler

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecases"
	"Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/templates"
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
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}

	template.AuthorizedChecker(ctx, reqData.CustomerId)
	newShopingCart, err := handler.Usecase.Add(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(newShopingCart))
}

func (handler *ShoppingCartHandler) GetListShoppingCart(ctx *gin.Context) {
	customerId, _ := strconv.Atoi(ctx.Query("customerId"))

	template.AuthorizedChecker(ctx, ctx.Query("customerId"))
	if ctx.IsAborted() {
		return
	}

	listShoppingCart, err := handler.Usecase.GetListItem(ctx, customerId)
	if err != nil {
		if err.Error() == "cart is empty" {
			ctx.JSON(http.StatusOK, response.Success("your shopping cart is empty"))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, response.Success(listShoppingCart))
}

func (handler *ShoppingCartHandler) DeleteShoppingList(ctx *gin.Context) {
	var data model.DeleteReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
		return
	}
	template.AuthorizedChecker(ctx, data.CustomerId)
	if ctx.IsAborted() {
		return
	}

	err := handler.Usecase.DeleteShoppingList(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Delete Shopping List"))
}
