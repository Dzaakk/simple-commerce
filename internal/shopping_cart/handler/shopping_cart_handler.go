package handler

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/template"
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
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	template.AuthorizedChecker(ctx, reqData.CustomerID)
	newShopingCart, err := handler.Usecase.Add(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(newShopingCart))
}
func (handler *ShoppingCartHandler) UpdateShoppingCart(ctx *gin.Context) {}

func (handler *ShoppingCartHandler) GetListShoppingCart(ctx *gin.Context) {
	customerID, _ := strconv.Atoi(ctx.Query("customerId"))

	template.AuthorizedChecker(ctx, ctx.Query("customerId"))
	if ctx.IsAborted() {
		return
	}

	listShoppingCart, err := handler.Usecase.GetListItem(ctx, customerID)
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
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}
	template.AuthorizedChecker(ctx, data.CustomerID)
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
