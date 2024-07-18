package shopping_cart

import (
	model "Dzaakk/synapsis/internal/shopping_cart/models"
	usecase "Dzaakk/synapsis/internal/shopping_cart/usecase"
	auth "Dzaakk/synapsis/package/auth"
	template "Dzaakk/synapsis/package/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShoppingCarthandler struct {
	Usecase usecase.ShoppingCartUseCase
}

func NewShoppingCartHandler(usecase usecase.ShoppingCartUseCase) *ShoppingCarthandler {
	return &ShoppingCarthandler{
		Usecase: usecase,
	}
}
func (handler *ShoppingCarthandler) Route(r *gin.RouterGroup) {
	ShoppingHandler := r.Group("api/v1")

	ShoppingHandler.Use()
	{
		ShoppingHandler.POST("/shopping-cart", auth.JWTMiddleware(), func(ctx *gin.Context) {
			if err := handler.AddProductToShoppingCart(ctx); err != nil {
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
