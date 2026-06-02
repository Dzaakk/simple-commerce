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

var customerColumns = []string{"id", "email", "password_hash", "full_name", "phone", "status", "created_at", "updated_at"}

func TestCustomerRepositoryCreate(t *testing.T) {
	now := time.Date(2026, time.June, 2, 10, 0, 0, 0, time.UTC)
	customer := &model.Customer{
		Email:        "customer@example.com",
		PasswordHash: "hashed-password",
		FullName:     "Customer Name",
		Phone:        "08123456789",
		Status:       string(constant.StatusPending),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	db := newTestDB(t, expectQuery(
		customerQueryCreate,
		[]any{customer.Email, customer.PasswordHash, customer.FullName, customer.Phone, customer.Status, customer.CreatedAt, customer.UpdatedAt},
		rows([]string{"id"}, []driver.Value{"customer-1"}),
	))

	got, err := NewCustomerRepository(db).Create(context.Background(), customer)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != "customer-1" {
		t.Fatalf("id = %q, want customer-1", got)
	}
}

func TestCustomerRepositoryCreateReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("insert failed")
	db := newTestDB(t, expectQueryError(
		customerQueryCreate,
		[]any{"customer@example.com", "hash", "Customer Name", "08123456789", string(constant.StatusPending), time.Time{}, time.Time{}},
		wantErr,
	))

	got, err := NewCustomerRepository(db).Create(context.Background(), &model.Customer{
		Email:        "customer@example.com",
		PasswordHash: "hash",
		FullName:     "Customer Name",
		Phone:        "08123456789",
		Status:       string(constant.StatusPending),
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != "" {
		t.Fatalf("id = %q, want empty", got)
	}
}

func TestCustomerRepositoryUpdate(t *testing.T) {
	now := time.Date(2026, time.June, 2, 11, 0, 0, 0, time.UTC)
	customer := &model.Customer{
		ID:        "customer-1",
		Email:     "new@example.com",
		FullName:  "New Name",
		Phone:     "08999999999",
		Status:    string(constant.StatusActive),
		UpdatedAt: now,
	}
	db := newTestDB(t, expectExec(
		customerQueryUpdate,
		[]any{customer.Email, customer.FullName, customer.Phone, customer.Status, customer.UpdatedAt, customer.ID},
		1,
	))

	got, err := NewCustomerRepository(db).Update(context.Background(), customer)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got != 1 {
		t.Fatalf("rows affected = %d, want 1", got)
	}
}

func TestCustomerRepositoryUpdateReturnsZeroWhenNoRowsAffected(t *testing.T) {
	db := newTestDB(t, expectExec(
		customerQueryUpdate,
		[]any{"new@example.com", "New Name", "08999999999", string(constant.StatusActive), time.Time{}, "missing"},
		0,
	))

	got, err := NewCustomerRepository(db).Update(context.Background(), &model.Customer{
		ID:       "missing",
		Email:    "new@example.com",
		FullName: "New Name",
		Phone:    "08999999999",
		Status:   string(constant.StatusActive),
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got != 0 {
		t.Fatalf("rows affected = %d, want 0", got)
	}
}

func TestCustomerRepositoryFindByID(t *testing.T) {
	now := time.Date(2026, time.June, 2, 12, 0, 0, 0, time.UTC)
	db := newTestDB(t, expectQuery(
		customerQueryFindByID,
		[]any{"customer-1"},
		rows(customerColumns, customerRow("customer-1", "customer@example.com", "Customer Name", "08123456789", string(constant.StatusActive), now)),
	))

	got, err := NewCustomerRepository(db).FindByID(context.Background(), "customer-1")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	assertCustomer(t, got, "customer-1", "customer@example.com", "Customer Name", "08123456789", string(constant.StatusActive), now)
}

func TestCustomerRepositoryFindByEmailReturnsNilWhenNotFound(t *testing.T) {
	db := newTestDB(t, expectQuery(
		customerQueryFindByEmail,
		[]any{"missing@example.com"},
		rows(customerColumns),
	))

	got, err := NewCustomerRepository(db).FindByEmail(context.Background(), "missing@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("customer = %#v, want nil", got)
	}
}

func TestCustomerRepositoryUpdateStatus(t *testing.T) {
	db := newTestDB(t, expectExec(
		customerQueryUpdateStatus,
		[]any{string(constant.StatusActive), "customer-1"},
		1,
	))

	if err := NewCustomerRepository(db).UpdateStatus(context.Background(), "customer-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
}

func TestCustomerRepositoryUpdateStatusReturnsErrorWhenNoRowsAffected(t *testing.T) {
	db := newTestDB(t, expectExec(
		customerQueryUpdateStatus,
		[]any{string(constant.StatusActive), "missing"},
		0,
	))

	err := NewCustomerRepository(db).UpdateStatus(context.Background(), "missing", constant.StatusActive)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
}

func TestCustomerRepositoryUpdateStatusWithTx(t *testing.T) {
	db := newTestDB(t, expectExec(
		customerQueryUpdateStatus,
		[]any{string(constant.StatusActive), "customer-1"},
		1,
	))
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	if err := NewCustomerRepository(db).UpdateStatusWithTx(context.Background(), tx, "customer-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatusWithTx returned error: %v", err)
	}
}

func customerRow(id, email, fullName, phone, status string, at time.Time) []driver.Value {
	return []driver.Value{id, email, "hashed-password", fullName, phone, status, at, at}
}

func assertCustomer(t *testing.T, got *model.Customer, id, email, fullName, phone, status string, at time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("customer is nil")
	}
	if got.ID != id || got.Email != email || got.FullName != fullName || got.Phone != phone ||
		got.Status != status || !got.CreatedAt.Equal(at) || !got.UpdatedAt.Equal(at) {
		t.Fatalf("customer = %#v", got)
	}
}
