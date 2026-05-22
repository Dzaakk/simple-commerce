package service

import (
	"testing"

	"Dzaakk/simple-commerce/internal/catalog/dto"
)

func TestProductListCacheKey(t *testing.T) {
	categoryID := int64(7)

	tests := []struct {
		name      string
		req       dto.ProductQueryReq
		wantKey   string
		cacheable bool
	}{
		{
			name: "general product list benchmark query",
			req: dto.ProductQueryReq{
				Limit: 100,
			},
			wantKey:   "catalog:v2:products:list:limit=100",
			cacheable: true,
		},
		{
			name: "category product list benchmark query",
			req: dto.ProductQueryReq{
				CategoryID: &categoryID,
				Limit:      100,
			},
			wantKey:   "catalog:v2:products:list:category_id=7:limit=100",
			cacheable: true,
		},
		{
			name: "cursor query is not cached",
			req: dto.ProductQueryReq{
				Cursor: stringPtr("2026-01-01T00:00:00Z|product-id"),
				Limit:  100,
			},
			cacheable: false,
		},
		{
			name: "non benchmark limit is not cached",
			req: dto.ProductQueryReq{
				Limit: 50,
			},
			cacheable: false,
		},
		{
			name: "additional filters are not cached",
			req: dto.ProductQueryReq{
				Name:  stringPtr("phone"),
				Limit: 100,
			},
			cacheable: false,
		},
		{
			name: "explicit sort is not cached",
			req: dto.ProductQueryReq{
				Limit:  100,
				SortBy: "price_asc",
			},
			cacheable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotCacheable := productListCacheKey(tt.req)
			if gotCacheable != tt.cacheable {
				t.Fatalf("cacheable = %v, want %v", gotCacheable, tt.cacheable)
			}
			if gotKey != tt.wantKey {
				t.Fatalf("key = %q, want %q", gotKey, tt.wantKey)
			}
		})
	}
}

func TestProductDetailCacheKey(t *testing.T) {
	productID := "6f52a63f-36c3-4f72-a1fb-8e384b490c6a"
	want := "catalog:v2:product:id:" + productID

	if got := productDetailCacheKey(productID); got != want {
		t.Fatalf("key = %q, want %q", got, want)
	}
}

func stringPtr(value string) *string {
	return &value
}
