package handler

import (
	"net/http"
	"strconv"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	products   service.ProductService
	categories service.CategoryService
}

func NewCatalogHandler(products service.ProductService, categories service.CategoryService) *CatalogHandler {
	return &CatalogHandler{products: products, categories: categories}
}

func (h *CatalogHandler) FindAllProducts(ctx *gin.Context) {
	h.findAllProducts(ctx, false)
}

func (h *CatalogHandler) FindAllProductsV2(ctx *gin.Context) {
	h.findAllProducts(ctx, true)
}

func (h *CatalogHandler) findAllProducts(ctx *gin.Context, cached bool) {
	query, err := productQuery(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var result *dto.ProductListResponse
	if cached {
		result, err = h.products.FindAllCached(ctx.Request.Context(), query)
	} else {
		result, err = h.products.FindAll(ctx.Request.Context(), query)
	}
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(result))
}

func (h *CatalogHandler) FindProductByID(ctx *gin.Context) {
	h.findProductByID(ctx, false)
}

func (h *CatalogHandler) FindProductByIDV2(ctx *gin.Context) {
	h.findProductByID(ctx, true)
}

func (h *CatalogHandler) findProductByID(ctx *gin.Context, cached bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid product id"))
		return
	}
	var result *dto.ProductResponse
	if cached {
		result, err = h.products.FindByIDCached(ctx.Request.Context(), id)
	} else {
		result, err = h.products.FindByID(ctx.Request.Context(), id)
	}
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(result))
}

func (h *CatalogHandler) FindAllCategories(ctx *gin.Context) {
	result, err := h.categories.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Success(result))
}
