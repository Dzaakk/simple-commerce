package handlers

import (
	usecase "Dzaakk/simple-commerce/internal/product/usecases"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Usecase usecase.ProductUseCase
}

func NewProductHandler(usecase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		Usecase: usecase,
	}
}

func (handler *ProductHandler) GetProduct(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Query("categoryId"))
	productName := ctx.Query("productName")

	if categoryId != 0 {
		listProduct, err := handler.Usecase.FindByCategoryId(ctx, categoryId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
			return
		}
		if listProduct == nil {
			ctx.JSON(http.StatusOK, response.Success("Product Not Found"))
			return
		}
		ctx.JSON(http.StatusOK, response.Success(listProduct))
	}
	if productName != "" {
		product, err := handler.Usecase.FindByName(ctx, productName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
			return
		} else if product == nil {
			ctx.JSON(http.StatusOK, response.Success("Product Not Found"))
			return
		}
		ctx.JSON(http.StatusOK, response.Success(product))
	}

}
