package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/model"
	"github.com/DATA-DOG/go-sqlmock"
)

var productColumns = []string{"id", "seller_id", "category_id", "name", "sku", "description", "price", "image_url", "is_active", "created_at", "updated_at"}

func TestProductRepositoryCreate(t *testing.T) {
	description := "Fast phone"
	imageURL := "https://example.com/phone.jpg"
	now := time.Date(2026, time.June, 3, 13, 0, 0, 0, time.UTC)
	product := &model.Product{
		SellerID:    "seller-1",
		CategoryID:  42,
		Name:        "Phone",
		SKU:         "PHONE-1",
		Description: &description,
		Price:       1500000,
		ImageURL:    &imageURL,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	db, mock := newMockDB(t)
	mock.ExpectQuery(productQueryCreate).
		WithArgs(product.SellerID, product.CategoryID, product.Name, product.SKU, description, product.Price, imageURL, product.IsActive, product.CreatedAt, product.UpdatedAt).
		WillReturnRows(sqlmockRows([]string{"id"}).AddRow("product-1"))

	got, err := NewProductRepository(db).Create(context.Background(), product)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != "product-1" {
		t.Fatalf("id = %q, want product-1", got)
	}
}

func TestProductRepositoryUpdateReturnsNoRowsError(t *testing.T) {
	now := time.Date(2026, time.June, 3, 14, 0, 0, 0, time.UTC)
	product := &model.Product{
		ID:         "missing",
		SellerID:   "seller-1",
		CategoryID: 42,
		Name:       "Phone",
		SKU:        "PHONE-1",
		Price:      1500000,
		IsActive:   true,
		UpdatedAt:  now,
	}
	db, mock := newMockDB(t)
	mock.ExpectExec(productQueryUpdate).
		WithArgs(product.SellerID, product.CategoryID, product.Name, product.SKU, nil, product.Price, nil, product.IsActive, product.UpdatedAt, product.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	got, err := NewProductRepository(db).Update(context.Background(), product)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
	if got != 0 {
		t.Fatalf("rows affected = %d, want 0", got)
	}
}

func TestProductRepositorySoftDelete(t *testing.T) {
	now := time.Date(2026, time.June, 3, 15, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectExec(productQuerySoftDelete).
		WithArgs(now, "product-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	got, err := NewProductRepository(db).SoftDelete(context.Background(), "product-1", now)
	if err != nil {
		t.Fatalf("SoftDelete returned error: %v", err)
	}
	if got != 1 {
		t.Fatalf("rows affected = %d, want 1", got)
	}
}

func TestProductRepositoryFindByIDReturnsNotFoundError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery(productQueryFindByID).
		WithArgs("missing").
		WillReturnRows(sqlmockRows(productColumns))

	got, err := NewProductRepository(db).FindByID(context.Background(), "missing")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
	if got != nil {
		t.Fatalf("product = %#v, want nil", got)
	}
}

func TestProductRepositoryFindBySellerID(t *testing.T) {
	description := "Fast phone"
	imageURL := "https://example.com/phone.jpg"
	now := time.Date(2026, time.June, 3, 16, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(productQueryFindBySeller).
		WithArgs("seller-1").
		WillReturnRows(sqlmockRows(productColumns).
			AddRow(productRow("product-1", "seller-1", 42, "Phone", "PHONE-1", &description, 1500000, &imageURL, true, now)...).
			AddRow(productRow("product-2", "seller-1", 43, "Case", "CASE-1", nil, 125000, nil, true, now)...))

	got, err := NewProductRepository(db).FindBySellerID(context.Background(), "seller-1")
	if err != nil {
		t.Fatalf("FindBySellerID returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("product count = %d, want 2", len(got))
	}
	assertProduct(t, got[0], "product-1", "seller-1", 42, "Phone", "PHONE-1", &description, 1500000, &imageURL, true, now)
	assertProduct(t, got[1], "product-2", "seller-1", 43, "Case", "CASE-1", nil, 125000, nil, true, now)
}

func TestProductRepositoryFindAllUsesBuiltQuery(t *testing.T) {
	categoryID := int64(42)
	sellerID := "seller-1"
	minPrice := 100000.0
	maxPrice := 2000000.0
	name := "phone"
	cursorTime := time.Date(2026, time.June, 3, 17, 0, 0, 0, time.UTC)
	cursor := cursorTime.Format(time.RFC3339Nano) + "|product-9"
	filter := ProductFilter{
		CategoryID: &categoryID,
		SellerID:   &sellerID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
		Name:       &name,
		Cursor:     &cursor,
		Limit:      20,
		SortBy:     "newest",
	}
	query, args := buildProductQuery(filter)
	description := "Fast phone"
	now := time.Date(2026, time.June, 3, 18, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(query).
		WithArgs(sqlmockArgs(args)...).
		WillReturnRows(sqlmockRows(productColumns).
			AddRow(productRow("product-1", "seller-1", 42, "Phone", "PHONE-1", &description, 1500000, nil, true, now)...))

	got, err := NewProductRepository(db).FindAll(context.Background(), filter)
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("product count = %d, want 1", len(got))
	}
	assertProduct(t, got[0], "product-1", "seller-1", 42, "Phone", "PHONE-1", &description, 1500000, nil, true, now)
}

func TestProductRepositoryUpdateStock(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(productQueryUpdateStock).
		WithArgs(int64(7), sqlmock.AnyArg(), "product-1", "seller-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := NewProductRepository(db).UpdateStock(context.Background(), "product-1", "seller-1", 7); err != nil {
		t.Fatalf("UpdateStock returned error: %v", err)
	}
}

func TestProductRepositoryUpdateStockReturnsNoRowsError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(productQueryUpdateStock).
		WithArgs(int64(7), sqlmock.AnyArg(), "product-1", "seller-1").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewProductRepository(db).UpdateStock(context.Background(), "product-1", "seller-1", 7)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func TestBuildProductQuery(t *testing.T) {
	categoryID := int64(42)
	sellerID := "seller-1"
	minPrice := 100000.0
	maxPrice := 2000000.0
	name := "phone"
	cursor := "1500000|product-9"

	gotQuery, gotArgs := buildProductQuery(ProductFilter{
		CategoryID: &categoryID,
		SellerID:   &sellerID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
		Name:       &name,
		Cursor:     &cursor,
		Limit:      20,
		SortBy:     "price_desc",
	})
	wantQuery := "SELECT " + productSelectColumns + " FROM public.products WHERE is_active = true" +
		" AND category_id = $1 AND seller_id = $2 AND price >= $3 AND price <= $4 AND name ILIKE $5" +
		" AND (price, id) < ($6, $7) ORDER BY price DESC, id DESC LIMIT $8"
	wantArgs := []any{categoryID, sellerID, minPrice, maxPrice, "%phone%", 1500000.0, "product-9", 20}

	if gotQuery != wantQuery {
		t.Fatalf("query = %q, want %q", gotQuery, wantQuery)
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func productRow(id, sellerID string, categoryID int64, name, sku string, description *string, price float64, imageURL *string, isActive bool, at time.Time) []driver.Value {
	var desc any
	if description != nil {
		desc = *description
	}
	var image any
	if imageURL != nil {
		image = *imageURL
	}
	return []driver.Value{id, sellerID, categoryID, name, sku, desc, price, image, isActive, at, at}
}

func assertProduct(t *testing.T, got *model.Product, id, sellerID string, categoryID int64, name, sku string, description *string, price float64, imageURL *string, isActive bool, at time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("product is nil")
	}
	if got.ID != id || got.SellerID != sellerID || got.CategoryID != categoryID ||
		got.Name != name || got.SKU != sku || got.Price != price || got.IsActive != isActive ||
		!got.CreatedAt.Equal(at) || !got.UpdatedAt.Equal(at) {
		t.Fatalf("product = %#v", got)
	}
	assertStringPtr(t, "description", got.Description, description)
	assertStringPtr(t, "image url", got.ImageURL, imageURL)
}

func assertStringPtr(t *testing.T, field string, got *string, want *string) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("%s = %q, want nil", field, *got)
		}
		return
	}
	if got == nil || *got != *want {
		t.Fatalf("%s = %v, want %q", field, got, *want)
	}
}
