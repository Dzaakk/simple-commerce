package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/package/response"

	"github.com/go-redis/redis/v8"
)

type ProductServiceImpl struct {
	repo  ProductRepository
	redis *redis.Client
}

func NewProductService(repo ProductRepository, redisClient *redis.Client) *ProductServiceImpl {
	return &ProductServiceImpl{repo: repo, redis: redisClient}
}

func (s *ProductServiceImpl) FindByID(ctx context.Context, id int64) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, response.NewAppError(http.StatusNotFound, "product not found")
	}
	result := dto.ToProductResponse(product)
	return &result, nil
}

func (s *ProductServiceImpl) FindByIDCached(ctx context.Context, id int64) (*dto.ProductResponse, error) {
	key := "catalog:product:" + strconv.FormatInt(id, 10)
	var result dto.ProductResponse
	if readCache(ctx, s.redis, key, &result) {
		return &result, nil
	}
	product, err := s.FindByID(ctx, id)
	if err == nil {
		writeCache(ctx, s.redis, key, product)
	}
	return product, err
}

func (s *ProductServiceImpl) FindAll(ctx context.Context, query dto.ProductQuery) (*dto.ProductListResponse, error) {
	products, err := s.repo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		items = append(items, dto.ToProductResponse(product))
	}
	result := &dto.ProductListResponse{Items: items}
	if len(items) == query.Limit && len(items) > 0 {
		last := items[len(items)-1]
		var cursor string
		switch query.SortBy {
		case "price_asc", "price_desc":
			cursor = strconv.FormatFloat(last.Price, 'f', -1, 64) + "|" + strconv.FormatInt(last.ID, 10)
		default:
			cursor = last.CreatedAt.Format(time.RFC3339Nano) + "|" + strconv.FormatInt(last.ID, 10)
		}
		result.NextCursor = &cursor
	}
	return result, nil
}

func (s *ProductServiceImpl) FindAllCached(ctx context.Context, query dto.ProductQuery) (*dto.ProductListResponse, error) {
	key := productListCacheKey(query)
	var result dto.ProductListResponse
	if readCache(ctx, s.redis, key, &result) {
		return &result, nil
	}
	products, err := s.FindAll(ctx, query)
	if err == nil {
		writeCache(ctx, s.redis, key, products)
	}
	return products, err
}
