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
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	id, err := h.ProductService.Create(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(id))
}

func (h *CatalogHandler) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	sellerID := ctx.Query("seller_id")
	if sellerID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	var req dto.UpdateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.ProductService.Update(ctx, productID, sellerID, &req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update Product"))
}

func (h *CatalogHandler) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	sellerID := ctx.Query("seller_id")
	if sellerID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.ProductService.SoftDelete(ctx, productID, sellerID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Delete Product"))
}

func (h *CatalogHandler) FindProductByID(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.ProductService.FindByID(ctx, productID)
	if err != nil {
		ctx.Error(err)
		return
	}
	if data == nil {
		ctx.Error(response.NewAppError(http.StatusNotFound, "product not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindAllProducts(ctx *gin.Context) {
	req, err := productQueryReqFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.ProductService.FindAll(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) UpdateProductStock(ctx *gin.Context) {
	req := dto.UpdateStockReq{
		ProductID: ctx.Param("id"),
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if req.ProductID == "" || req.SellerID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	if err := h.ProductService.UpdateStock(ctx, req.ProductID, req.SellerID, req.Quantity); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Success Update Product Stock"))
}

func (h *CatalogHandler) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	id, err := h.CategoryService.Create(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(id))
}

func (h *CatalogHandler) FindAllCategories(ctx *gin.Context) {
	data, err := h.CategoryService.FindAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindCategoryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.CategoryService.FindByID(ctx, id)
	if err != nil {
		ctx.Error(err)
		return
	}
	if data == nil {
		ctx.Error(response.NewAppError(http.StatusNotFound, "category not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
