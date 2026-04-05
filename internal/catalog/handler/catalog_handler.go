package handler

import (
	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	ProductService  service.ProductService
	CategoryService service.CategoryService
}

func NewCatalogHandler(productService service.ProductService, categoryService service.CategoryService) *CatalogHandler {
	return &CatalogHandler{
		ProductService:  productService,
		CategoryService: categoryService,
	}
}

func (h *CatalogHandler) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	id, err := h.ProductService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(id))
}

func (h *CatalogHandler) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	sellerID := ctx.Query("seller_id")
	if sellerID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	var req dto.UpdateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.ProductService.Update(ctx, productID, sellerID, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update Product"))
}

func (h *CatalogHandler) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	sellerID := ctx.Query("seller_id")
	if sellerID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.ProductService.SoftDelete(ctx, productID, sellerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Delete Product"))
}

func (h *CatalogHandler) FindProductByID(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.ProductService.FindByID(ctx, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NotFound("product not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindAllProducts(ctx *gin.Context) {
	var req dto.ProductQueryReq

	if val := ctx.Query("category_id"); val != "" {
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
			return
		}
		req.CategoryID = &id
	}
	if val := ctx.Query("seller_id"); val != "" {
		req.SellerID = &val
	}
	if val := ctx.Query("min_price"); val != "" {
		minPrice, err := strconv.ParseFloat(val, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
			return
		}
		req.MinPrice = &minPrice
	}
	if val := ctx.Query("max_price"); val != "" {
		maxPrice, err := strconv.ParseFloat(val, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
			return
		}
		req.MaxPrice = &maxPrice
	}
	if val := ctx.Query("name"); val != "" {
		req.Name = &val
	}
	if val := ctx.Query("cursor"); val != "" {
		req.Cursor = &val
	}
	if val := ctx.Query("limit"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
			return
		}
		req.Limit = limit
	}
	req.SortBy = ctx.Query("sort_by")

	data, err := h.ProductService.FindAll(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

type updateStockReq struct {
	Quantity int `json:"quantity"`
}

func (h *CatalogHandler) UpdateProductStock(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	sellerID := ctx.Query("seller_id")
	if sellerID == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	var req updateStockReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	if err := h.ProductService.UpdateStock(ctx, productID, sellerID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update Product Stock"))
}

func (h *CatalogHandler) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	id, err := h.CategoryService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(id))
}

func (h *CatalogHandler) FindAllCategories(ctx *gin.Context) {
	data, err := h.CategoryService.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindCategoryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.CategoryService.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NotFound("category not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
