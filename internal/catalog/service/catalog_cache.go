package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"

	"github.com/go-redis/redis/v8"
)

const catalogCacheTTL = 5 * time.Minute

func readCache(ctx context.Context, client *redis.Client, key string, target any) bool {
	if client == nil {
		return false
	}
	value, err := client.Get(ctx, key).Bytes()
	return err == nil && json.Unmarshal(value, target) == nil
}

func writeCache(ctx context.Context, client *redis.Client, key string, value any) {
	if client == nil {
		return
	}
	payload, err := json.Marshal(value)
	if err == nil {
		_ = client.Set(ctx, key, payload, catalogCacheTTL).Err()
	}
}

func productListCacheKey(query dto.ProductQuery) string {
	payload, _ := json.Marshal(query)
	return fmt.Sprintf("catalog:products:%x", payload)
}
