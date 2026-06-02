package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
)

var sellerColumns = []string{"id", "email", "password_hash", "shop_name", "phone", "status", "created_at", "updated_at"}

func TestSellerRepositoryCreate(t *testing.T) {
	now := time.Date(2026, time.June, 2, 10, 0, 0, 0, time.UTC)
	seller := &model.Seller{
		Email:        "seller@example.com",
		PasswordHash: "hashed-password",
		ShopName:     "Best Shop",
		Phone:        "08123456789",
		Status:       string(constant.StatusPending),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	db := newTestDB(t, expectQuery(
		sellerQueryCreate,
		[]any{seller.Email, seller.PasswordHash, seller.ShopName, seller.Phone, seller.Status, seller.CreatedAt, seller.UpdatedAt},
		rows([]string{"id"}, []driver.Value{"seller-1"}),
	))

	got, err := NewSellerRepository(db).Create(context.Background(), seller)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != "seller-1" {
		t.Fatalf("id = %q, want seller-1", got)
	}
}

func TestSellerRepositoryUpdate(t *testing.T) {
	now := time.Date(2026, time.June, 2, 11, 0, 0, 0, time.UTC)
	seller := &model.Seller{
		ID:        "seller-1",
		Email:     "new@example.com",
		ShopName:  "New Shop",
		Phone:     "08999999999",
		Status:    string(constant.StatusActive),
		UpdatedAt: now,
	}
	db := newTestDB(t, expectExec(
		sellerQueryUpdate,
		[]any{seller.Email, seller.ShopName, seller.Phone, seller.Status, seller.UpdatedAt, seller.ID},
		1,
	))

	got, err := NewSellerRepository(db).Update(context.Background(), seller)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got != 1 {
		t.Fatalf("rows affected = %d, want 1", got)
	}
}

func TestSellerRepositoryUpdateReturnsWrappedExecError(t *testing.T) {
	wantErr := errors.New("update failed")
	db := newTestDB(t, expectExecError(
		sellerQueryUpdate,
		[]any{"new@example.com", "New Shop", "08999999999", string(constant.StatusActive), time.Time{}, "seller-1"},
		wantErr,
	))

	got, err := NewSellerRepository(db).Update(context.Background(), &model.Seller{
		ID:       "seller-1",
		Email:    "new@example.com",
		ShopName: "New Shop",
		Phone:    "08999999999",
		Status:   string(constant.StatusActive),
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("rows affected = %d, want 0", got)
	}
}

func TestSellerRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	db := newTestDB(t, expectQuery(
		sellerQueryFindByID,
		[]any{"missing"},
		rows(sellerColumns),
	))

	got, err := NewSellerRepository(db).FindByID(context.Background(), "missing")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("seller = %#v, want nil", got)
	}
}

func TestSellerRepositoryFindByEmail(t *testing.T) {
	now := time.Date(2026, time.June, 2, 12, 0, 0, 0, time.UTC)
	db := newTestDB(t, expectQuery(
		sellerQueryFindByEmail,
		[]any{"seller@example.com"},
		rows(sellerColumns, sellerRow("seller-1", "seller@example.com", "Best Shop", "08123456789", string(constant.StatusActive), now)),
	))

	got, err := NewSellerRepository(db).FindByEmail(context.Background(), "seller@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	assertSeller(t, got, "seller-1", "seller@example.com", "Best Shop", "08123456789", string(constant.StatusActive), now)
}

func TestSellerRepositoryFindByShopName(t *testing.T) {
	now := time.Date(2026, time.June, 2, 13, 0, 0, 0, time.UTC)
	db := newTestDB(t, expectQuery(
		sellerQueryFindByName,
		[]any{"%shop%"},
		rows(
			sellerColumns,
			sellerRow("seller-1", "first@example.com", "First Shop", "08111111111", string(constant.StatusActive), now),
			sellerRow("seller-2", "second@example.com", "Second Shop", "08222222222", string(constant.StatusPending), now),
		),
	))

	got, err := NewSellerRepository(db).FindByShopName(context.Background(), "shop")
	if err != nil {
		t.Fatalf("FindByShopName returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("seller count = %d, want 2", len(got))
	}
	assertSeller(t, got[0], "seller-1", "first@example.com", "First Shop", "08111111111", string(constant.StatusActive), now)
	assertSeller(t, got[1], "seller-2", "second@example.com", "Second Shop", "08222222222", string(constant.StatusPending), now)
}

func TestSellerRepositoryFindByShopNameReturnsQueryError(t *testing.T) {
	wantErr := errors.New("select failed")
	db := newTestDB(t, expectQueryError(sellerQueryFindByName, []any{"%shop%"}, wantErr))

	got, err := NewSellerRepository(db).FindByShopName(context.Background(), "shop")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != nil {
		t.Fatalf("sellers = %#v, want nil", got)
	}
}

func TestSellerRepositoryUpdateStatusReturnsErrorWhenNoRowsAffected(t *testing.T) {
	db := newTestDB(t, expectExec(
		sellerQueryUpdateStatus,
		[]any{string(constant.StatusActive), "missing"},
		0,
	))

	err := NewSellerRepository(db).UpdateStatus(context.Background(), "missing", constant.StatusActive)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func TestSellerRepositoryUpdateStatusWithTx(t *testing.T) {
	db := newTestDB(t, expectExec(
		sellerQueryUpdateStatus,
		[]any{string(constant.StatusActive), "seller-1"},
		1,
	))
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	if err := NewSellerRepository(db).UpdateStatusWithTx(context.Background(), tx, "seller-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatusWithTx returned error: %v", err)
	}
}

func sellerRow(id, email, shopName, phone, status string, at time.Time) []driver.Value {
	return []driver.Value{id, email, "hashed-password", shopName, phone, status, at, at}
}

func assertSeller(t *testing.T, got *model.Seller, id, email, shopName, phone, status string, at time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("seller is nil")
	}
	if got.ID != id || got.Email != email || got.ShopName != shopName || got.Phone != phone ||
		got.Status != status || !got.CreatedAt.Equal(at) || !got.UpdatedAt.Equal(at) {
		t.Fatalf("seller = %#v", got)
	}
}
