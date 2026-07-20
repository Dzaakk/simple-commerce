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
	"Dzaakk/simple-commerce/package/db/transactor"
	"github.com/DATA-DOG/go-sqlmock"
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
	db, mock := newMockDB(t)
	mock.ExpectQuery(sellerQueryCreate).
		WithArgs(seller.Email, seller.PasswordHash, seller.ShopName, seller.Phone, seller.Status, seller.CreatedAt, seller.UpdatedAt).
		WillReturnRows(sqlmockRows([]string{"id"}).AddRow("seller-1"))

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
	db, mock := newMockDB(t)
	mock.ExpectExec(sellerQueryUpdate).
		WithArgs(seller.Email, seller.ShopName, seller.Phone, seller.Status, seller.UpdatedAt, seller.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

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
	db, mock := newMockDB(t)
	mock.ExpectExec(sellerQueryUpdate).
		WithArgs("new@example.com", "New Shop", "08999999999", string(constant.StatusActive), time.Time{}, "seller-1").
		WillReturnError(wantErr)

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
	db, mock := newMockDB(t)
	mock.ExpectQuery(sellerQueryFindByID).
		WithArgs("missing").
		WillReturnRows(sqlmockRows(sellerColumns))

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
	db, mock := newMockDB(t)
	mock.ExpectQuery(sellerQueryFindByEmail).
		WithArgs("seller@example.com").
		WillReturnRows(sqlmockRows(sellerColumns).
			AddRow(sellerRow("seller-1", "seller@example.com", "Best Shop", "08123456789", string(constant.StatusActive), now)...))

	got, err := NewSellerRepository(db).FindByEmail(context.Background(), "seller@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	assertSeller(t, got, "seller-1", "seller@example.com", "Best Shop", "08123456789", string(constant.StatusActive), now)
}

func TestSellerRepositoryFindByShopName(t *testing.T) {
	now := time.Date(2026, time.June, 2, 13, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(sellerQueryFindByName).
		WithArgs("%shop%").
		WillReturnRows(sqlmockRows(sellerColumns).
			AddRow(sellerRow("seller-1", "first@example.com", "First Shop", "08111111111", string(constant.StatusActive), now)...).
			AddRow(sellerRow("seller-2", "second@example.com", "Second Shop", "08222222222", string(constant.StatusPending), now)...))

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
	db, mock := newMockDB(t)
	mock.ExpectQuery(sellerQueryFindByName).WithArgs("%shop%").WillReturnError(wantErr)

	got, err := NewSellerRepository(db).FindByShopName(context.Background(), "shop")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != nil {
		t.Fatalf("sellers = %#v, want nil", got)
	}
}

func TestSellerRepositoryUpdateStatusReturnsErrorWhenNoRowsAffected(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(sellerQueryUpdateStatus).
		WithArgs(string(constant.StatusActive), "missing").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewSellerRepository(db).UpdateStatus(context.Background(), "missing", constant.StatusActive)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func TestSellerRepositoryUpdateStatusUsesContextExecutor(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(sellerQueryUpdateStatus).
		WithArgs(string(constant.StatusActive), "seller-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	ctx := transactor.WithExecutor(context.Background(), tx)
	if err := NewSellerRepository(db).UpdateStatus(ctx, "seller-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
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
