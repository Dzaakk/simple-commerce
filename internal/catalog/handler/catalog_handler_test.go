package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

type mockProductService struct {
	createFn      func(context.Context, *dto.CreateProductReq) (string, error)
	updateFn      func(context.Context, string, string, *dto.UpdateProductReq) error
	softDeleteFn  func(context.Context, string, string) error
	findByIDFn    func(context.Context, string) (*dto.ProductRes, error)
	findAllFn     func(context.Context, dto.ProductQueryReq) (*dto.ProductListRes, error)
	updateStockFn func(context.Context, string, string, int) error
}

func (m *mockProductService) Create(ctx context.Context, req *dto.CreateProductReq) (string, error) {
	if m.createFn == nil {
		return "", errors.New("unexpected Create call")
	}
	return m.createFn(ctx, req)
}

func (m *mockProductService) Update(ctx context.Context, productID string, sellerID string, req *dto.UpdateProductReq) error {
	if m.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return m.updateFn(ctx, productID, sellerID, req)
}

func (m *mockProductService) SoftDelete(ctx context.Context, productID string, sellerID string) error {
	if m.softDeleteFn == nil {
		return errors.New("unexpected SoftDelete call")
	}
	return m.softDeleteFn(ctx, productID, sellerID)
}

func (m *mockProductService) FindByID(ctx context.Context, productID string) (*dto.ProductRes, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return m.findByIDFn(ctx, productID)
}

func (m *mockProductService) FindByIDCached(ctx context.Context, productID string) (*dto.ProductRes, error) {
	return m.FindByID(ctx, productID)
}

func (m *mockProductService) FindAll(ctx context.Context, req dto.ProductQueryReq) (*dto.ProductListRes, error) {
	if m.findAllFn == nil {
		return nil, errors.New("unexpected FindAll call")
	}
	return m.findAllFn(ctx, req)
}

func (m *mockProductService) FindAllCached(ctx context.Context, req dto.ProductQueryReq) (*dto.ProductListRes, error) {
	return m.FindAll(ctx, req)
}

func (m *mockProductService) UpdateStock(ctx context.Context, productID string, sellerID string, quantity int) error {
	if m.updateStockFn == nil {
		return errors.New("unexpected UpdateStock call")
	}
	return m.updateStockFn(ctx, productID, sellerID, quantity)
}

func TestCatalogHandlerUpdateProductUsesAuthenticatedSellerID(t *testing.T) {
	handler := NewCatalogHandler(&mockProductService{
		updateFn: func(_ context.Context, productID string, sellerID string, _ *dto.UpdateProductReq) error {
			if productID != "product-1" {
				t.Fatalf("product id = %q, want product-1", productID)
			}
			if sellerID != "seller-1" {
				t.Fatalf("seller id = %q, want seller-1", sellerID)
			}
			return nil
		},
	}, nil)

	w, ctx := performCatalogRequest(http.MethodPut, "/product/product-1", `{"name":"Phone"}`, handler.UpdateProduct)
	ctx.Params = gin.Params{{Key: "id", Value: "product-1"}}
	ctx.Set("id", "seller-1")

	handler.UpdateProduct(ctx)

	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCatalogHandlerUpdateProductRejectsMismatchedQuerySellerID(t *testing.T) {
	called := false
	handler := NewCatalogHandler(&mockProductService{
		updateFn: func(context.Context, string, string, *dto.UpdateProductReq) error {
			called = true
			return nil
		},
	}, nil)

	_, ctx := performCatalogRequest(http.MethodPut, "/product/product-1?seller_id=attacker", `{"name":"Phone"}`, handler.UpdateProduct)
	ctx.Params = gin.Params{{Key: "id", Value: "product-1"}}
	ctx.Set("id", "seller-1")

	handler.UpdateProduct(ctx)

	if called {
		t.Fatal("service Update must not be called for mismatched seller id")
	}
	assertCatalogHandlerAppError(t, ctx, http.StatusUnauthorized, "unauthorized")
}

func TestCatalogHandlerDeleteProductUsesAuthenticatedSellerID(t *testing.T) {
	handler := NewCatalogHandler(&mockProductService{
		softDeleteFn: func(_ context.Context, productID string, sellerID string) error {
			if productID != "product-1" {
				t.Fatalf("product id = %q, want product-1", productID)
			}
			if sellerID != "seller-1" {
				t.Fatalf("seller id = %q, want seller-1", sellerID)
			}
			return nil
		},
	}, nil)

	w, ctx := performCatalogRequest(http.MethodDelete, "/product/product-1", "", handler.DeleteProduct)
	ctx.Params = gin.Params{{Key: "id", Value: "product-1"}}
	ctx.Set("id", "seller-1")

	handler.DeleteProduct(ctx)

	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCatalogHandlerCreateProductRejectsMismatchedBodySellerID(t *testing.T) {
	called := false
	handler := NewCatalogHandler(&mockProductService{
		createFn: func(context.Context, *dto.CreateProductReq) (string, error) {
			called = true
			return "product-1", nil
		},
	}, nil)

	_, ctx := performCatalogRequest(http.MethodPost, "/product", `{"seller_id":"attacker","category_id":1,"name":"Phone","sku":"PHONE-1","price":1000}`, handler.CreateProduct)
	ctx.Set("id", "seller-1")

	handler.CreateProduct(ctx)

	if called {
		t.Fatal("service Create must not be called for mismatched seller id")
	}
	assertCatalogHandlerAppError(t, ctx, http.StatusUnauthorized, "unauthorized")
}

func performCatalogRequest(method, target, body string, _ gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	return w, ctx
}

func assertCatalogHandlerAppError(t *testing.T, ctx *gin.Context, code int, message string) {
	t.Helper()

	if len(ctx.Errors) != 1 {
		t.Fatalf("errors = %v, want exactly one", ctx.Errors)
	}
	var appErr *response.AppError
	if !errors.As(ctx.Errors[0].Err, &appErr) {
		t.Fatalf("error = %T %v, want *response.AppError", ctx.Errors[0].Err, ctx.Errors[0].Err)
	}
	if appErr.Code != code {
		t.Fatalf("code = %d, want %d", appErr.Code, code)
	}
	if appErr.Message != message {
		t.Fatalf("message = %q, want %q", appErr.Message, message)
	}
}
