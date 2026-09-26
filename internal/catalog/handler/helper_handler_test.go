package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductQueryDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/product", nil)
	query, err := productQuery(ctx)
	if err != nil || query.Limit != 20 || query.SortBy != "newest" {
		t.Fatalf("query = %#v, error = %v", query, err)
	}
}

func TestProductQueryRejectsInvalidPriceRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/product?min_price=20&max_price=10", nil)
	if _, err := productQuery(ctx); err == nil {
		t.Fatal("productQuery() error = nil")
	}
}

func TestProductQueryRejectsMalformedCursor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/product?cursor=not-a-cursor", nil)
	if _, err := productQuery(ctx); err == nil {
		t.Fatal("productQuery() error = nil")
	}
}
