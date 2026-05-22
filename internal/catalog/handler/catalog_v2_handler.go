package handler

import (
	"net/http"

	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

func (h *CatalogHandler) FindAllProductsV2(ctx *gin.Context) {
	req, err := productQueryReqFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.ProductService.FindAllCached(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindProductByIDV2(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.ProductService.FindByIDCached(ctx.Request.Context(), productID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindAllCategoriesV2(ctx *gin.Context) {
	data, err := h.CategoryService.FindAllCached(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
