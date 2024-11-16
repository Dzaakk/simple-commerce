package handler

import (
	usecase "Dzaakk/simple-commerce/internal/product/usecases"
	template "Dzaakk/simple-commerce/package/templates"
	"fmt"
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
		listProduct, err := handler.Usecase.FindByCategoryId(categoryId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		}
		if listProduct == nil {
			ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", fmt.Sprintf("product with category %v is not found", categoryId)))
			ctx.Abort()
		}
		ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success", listProduct))
	}
	if productName != "" {
		product, err := handler.Usecase.FindByName(productName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		} else if product == nil {
			ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Product Not Found", product))
		}
		ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success", product))
	}

}
