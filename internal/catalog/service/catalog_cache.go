package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"

	"github.com/go-redis/redis/v8"
)

const (
	catalogProductCacheTTL       = time.Minute
	catalogCategoryCacheTTL      = 5 * time.Minute
	catalogProductListCacheLimit = 100
)

func readCatalogCache(ctx context.Context, redisClient *redis.Client, key string, dst interface{}) bool {
	if redisClient == nil {
		return false
	}

	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}

	if err := json.Unmarshal(data, dst); err != nil {
		redisClient.Del(ctx, key)
		return false
	}

	return true
}

func writeCatalogCache(ctx context.Context, redisClient *redis.Client, key string, value interface{}, ttl time.Duration) {
	if redisClient == nil {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	redisClient.Set(ctx, key, data, ttl)
}

func productDetailCacheKey(productID string) string {
	return "catalog:v2:product:id:" + productID
}

func productListCacheKey(req dto.ProductQueryReq) (string, bool) {
	if req.Cursor != nil && *req.Cursor != "" {
		return "", false
	}
	if req.SellerID != nil || req.MinPrice != nil || req.MaxPrice != nil || req.Name != nil {
		return "", false
	}
	if req.SortBy != "" {
		return "", false
	}
	if req.Limit != catalogProductListCacheLimit {
		return "", false
	}

	if req.CategoryID != nil {
		return "catalog:v2:products:list:category_id=" +
			strconv.FormatInt(*req.CategoryID, 10) +
			":limit=" +
			strconv.Itoa(req.Limit), true
	}

	return "catalog:v2:products:list:limit=" + strconv.Itoa(req.Limit), true
}
