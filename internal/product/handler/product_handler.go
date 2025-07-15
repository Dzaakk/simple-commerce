package handler

import (
	"Dzaakk/simple-commerce/internal/product/model"
	"Dzaakk/simple-commerce/internal/product/usecase"
	"Dzaakk/simple-commerce/package/response"
	"net/http"

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

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	var params model.ProductFilter

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	listProduct, err := h.Usecase.FindByFilter(ctx, params)
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
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {}
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {}
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {}
