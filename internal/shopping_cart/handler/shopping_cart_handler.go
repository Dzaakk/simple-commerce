package shopping_cart

import (
	model "Dzaakk/synapsis/internal/shopping_cart/models"
	usecase "Dzaakk/synapsis/internal/shopping_cart/usecase"
	auth "Dzaakk/synapsis/package/auth"
	template "Dzaakk/synapsis/package/template"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ShoppingCarthandler struct {
	Usecase usecase.ShoppingCartUseCase
}

func NewShoppingCartHandler(usecase usecase.ShoppingCartUseCase) *ShoppingCarthandler {
	return &ShoppingCarthandler{
		Usecase: usecase,
	}
}
func (handler *ShoppingCarthandler) Route(r *gin.RouterGroup, redis *redis.Client) {
	ShoppingHandler := r.Group("api/v1")

	ShoppingHandler.Use()
	{
		ShoppingHandler.POST("/shopping-cart", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.AddProductToShoppingCart(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "internal server error", err.Error()))
				return
			}
		})
		ShoppingHandler.GET("/shopping-cart", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.GetListShoppingCart(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to Get product", err.Error()))
				return
			}
		})
		ShoppingHandler.POST("/shopping-cart/delete", auth.JWTMiddleware(redis), func(ctx *gin.Context) {
			if err := handler.DeleteShoppingList(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to Get product", err.Error()))
				return
			}
		})
	}
}

func (handler *ShoppingCarthandler) AddProductToShoppingCart(ctx *gin.Context) error {
	var reqData model.ShoppingCartReq
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return nil
	}

	template.AuthorizedChecker(ctx, reqData.CustomerId)
	newShopingCart, err := handler.Usecase.Add(reqData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return nil
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", newShopingCart))
	return nil
}

func (handler *ShoppingCarthandler) GetListShoppingCart(ctx *gin.Context) error {
	customerId, _ := strconv.Atoi(ctx.Query("customerId"))

	template.AuthorizedChecker(ctx, ctx.Query("customerId"))
	if ctx.IsAborted() {
		return nil
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
		return nil
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", listShoppingCart))
	return nil
}

func (handler *ShoppingCarthandler) DeleteShoppingList(ctx *gin.Context) error {
	var data model.DeleteReq
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, template.Response(http.StatusBadRequest, "Bad Request", "Invalid input data"))
		return nil
	}
	fmt.Printf("input = %v", data)
	template.AuthorizedChecker(ctx, data.CustomerId)
	if ctx.IsAborted() {
		return nil
	}

	err := handler.Usecase.DeleteShoppingList(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return nil
	}

	ctx.JSON(http.StatusCreated, template.Response(http.StatusOK, "Success", "Success Delete Shopping List"))
	return nil
}
