package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/db/transactor"
	"github.com/DATA-DOG/go-sqlmock"
)

var inventoryColumns = []string{"id", "product_id", "stock_quantity", "reserved_quantity", "version", "created_at", "updated_at"}

func TestInventoryRepositoryFindByProductID(t *testing.T) {
	now := time.Date(2026, time.June, 3, 12, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(inventoryQueryFindByProduct).
		WithArgs("product-1").
		WillReturnRows(sqlmockRows(inventoryColumns).AddRow(inventoryRow(1, "product-1", 20, 5, 2, now)...))

	got, err := NewInventoryRepository(db).FindByProductID(context.Background(), "product-1")
	if err != nil {
		t.Fatalf("FindByProductID returned error: %v", err)
	}
	assertInventory(t, got, 1, "product-1", 20, 5, 2, now)
}

func TestInventoryRepositoryFindByProductIDReturnsNilWhenNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery(inventoryQueryFindByProduct).
		WithArgs("missing").
		WillReturnRows(sqlmockRows(inventoryColumns))

	got, err := NewInventoryRepository(db).FindByProductID(context.Background(), "missing")
	if err != nil {
		t.Fatalf("FindByProductID returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("inventory = %#v, want nil", got)
	}
}

func TestInventoryRepositoryReserveStockRejectsInvalidInput(t *testing.T) {
	repo := NewInventoryRepository(nil)

	if err := repo.ReserveStock(context.Background(), "product-1", 0); err == nil || err.Error() != "invalid parameter quantity" {
		t.Fatalf("invalid qty error = %v, want invalid parameter quantity", err)
	}
}

func TestInventoryRepositoryReserveStock(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(inventoryQueryReserveStock).
		WithArgs(int64(3), sqlmock.AnyArg(), "product-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	ctx := transactor.WithExecutor(context.Background(), tx)
	if err := NewInventoryRepository(db).ReserveStock(ctx, "product-1", 3); err != nil {
		t.Fatalf("ReserveStock returned error: %v", err)
	}
}

func TestInventoryRepositoryReserveStockReturnsInsufficientStockWhenNoRowsAffected(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(inventoryQueryReserveStock).
		WithArgs(int64(3), sqlmock.AnyArg(), "product-1").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	ctx := transactor.WithExecutor(context.Background(), tx)
	err = NewInventoryRepository(db).ReserveStock(ctx, "product-1", 3)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func TestInventoryRepositoryReleaseStock(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(inventoryQueryReleaseStock).
		WithArgs(int64(2), sqlmock.AnyArg(), "product-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	ctx := transactor.WithExecutor(context.Background(), tx)
	if err := NewInventoryRepository(db).ReleaseStock(ctx, "product-1", 2); err != nil {
		t.Fatalf("ReleaseStock returned error: %v", err)
	}
}

func TestInventoryRepositoryReleaseStockReturnsNoRowsError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(inventoryQueryReleaseStock).
		WithArgs(int64(2), sqlmock.AnyArg(), "product-1").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	ctx := transactor.WithExecutor(context.Background(), tx)
	err = NewInventoryRepository(db).ReleaseStock(ctx, "product-1", 2)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func inventoryRow(id int64, productID string, stockQuantity, reservedQuantity, version int, at time.Time) []driver.Value {
	return []driver.Value{id, productID, int64(stockQuantity), int64(reservedQuantity), int64(version), at, at}
}

func assertInventory(t *testing.T, got *model.Inventory, id int64, productID string, stockQuantity, reservedQuantity, version int, at time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("inventory is nil")
	}
	if got.ID != id || got.ProductID != productID || got.StockQuantity != stockQuantity ||
		got.ReservedQuantity != reservedQuantity || got.Version != version ||
		!got.CreatedAt.Equal(at) || !got.UpdatedAt.Equal(at) {
		t.Fatalf("inventory = %#v", got)
	}
}
