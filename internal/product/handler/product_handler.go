package handler

import (
	"Dzaakk/simple-commerce/internal/product/usecase"
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
	categoryID, _ := strconv.Atoi(ctx.Query("categoryId"))
	productName := ctx.Query("productName")

	if categoryID != 0 {
		listProduct, err := handler.Usecase.FindByCategoryID(ctx, categoryID)
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
		product, err := handler.Usecase.FindByProductName(ctx, productName)
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
