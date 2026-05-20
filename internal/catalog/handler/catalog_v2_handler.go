package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	catalogProductCacheTTL  = time.Minute
	catalogCategoryCacheTTL = 5 * time.Minute
)

func (h *CatalogHandler) FindAllProductsV2(ctx *gin.Context) {
	req, err := productQueryReqFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	cacheKey := productListCacheKey(req)
	var cached dto.ProductListRes
	if h.getCatalogCache(ctx, cacheKey, &cached) {
		ctx.JSON(http.StatusOK, response.Success(&cached))
		return
	}

	data, err := h.ProductService.FindAll(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setCatalogCache(ctx, cacheKey, data, catalogProductCacheTTL)

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindProductByIDV2(ctx *gin.Context) {
	productID := ctx.Param("id")
	if productID == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	cacheKey := "catalog:v2:product:id:" + productID
	var cached dto.ProductRes
	if h.getCatalogCache(ctx, cacheKey, &cached) {
		ctx.JSON(http.StatusOK, response.Success(&cached))
		return
	}

	data, err := h.ProductService.FindByID(ctx, productID)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setCatalogCache(ctx, cacheKey, data, catalogProductCacheTTL)

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) FindAllCategoriesV2(ctx *gin.Context) {
	cacheKey := "catalog:v2:category:all"
	var cached []*dto.CategoryTree
	if h.getCatalogCache(ctx, cacheKey, &cached) {
		ctx.JSON(http.StatusOK, response.Success(cached))
		return
	}

	data, err := h.CategoryService.FindAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setCatalogCache(ctx, cacheKey, data, catalogCategoryCacheTTL)

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *CatalogHandler) getCatalogCache(ctx *gin.Context, key string, dst interface{}) bool {
	if h.Redis == nil {
		return false
	}

	requestCtx := ctx.Request.Context()

	data, err := h.Redis.Get(requestCtx, key).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return false
		}
		return false
	}

	if err := json.Unmarshal(data, dst); err != nil {
		h.Redis.Del(requestCtx, key)
		return false
	}

	return true
}

func (h *CatalogHandler) setCatalogCache(ctx *gin.Context, key string, value interface{}, ttl time.Duration) {
	if h.Redis == nil {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	h.Redis.Set(ctx.Request.Context(), key, data, ttl)
}

func productListCacheKey(req dto.ProductQueryReq) string {
	data, err := json.Marshal(req)
	if err != nil {
		return "catalog:v2:product:list:fallback"
	}

	sum := sha256.Sum256(data)

	return "catalog:v2:product:list:" + hex.EncodeToString(sum[:])
}
