package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/auth/model"
	"github.com/DATA-DOG/go-sqlmock"
)

var activationCodeColumns = []string{"id", "code", "email", "type", "user_type", "expires_at", "used_at", "created_at"}

func TestActivationCodeRepositoryCreate(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 4, 10, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.June, 4, 9, 0, 0, 0, time.UTC)
	code := &model.ActivationCode{
		Code:      "123456",
		Email:     "customer@example.com",
		Type:      "email_verification",
		UserType:  "customer",
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}
	db, mock := newMockDB(t)
	mock.ExpectQuery(activationCodeQueryCreate).
		WithArgs(code.Code, code.Email, code.Type, code.UserType, code.ExpiresAt, nil, code.CreatedAt).
		WillReturnRows(sqlmockRows([]string{"id"}).AddRow(int64(7)))

	got, err := NewActivationCodeRepository(db).Create(context.Background(), code)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != 7 {
		t.Fatalf("id = %d, want 7", got)
	}
}

func TestActivationCodeRepositoryCreateReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("insert failed")
	db, mock := newMockDB(t)
	mock.ExpectQuery(activationCodeQueryCreate).
		WithArgs("123456", "customer@example.com", "email_verification", "customer", time.Time{}, nil, time.Time{}).
		WillReturnError(wantErr)

	got, err := NewActivationCodeRepository(db).Create(context.Background(), &model.ActivationCode{
		Code:     "123456",
		Email:    "customer@example.com",
		Type:     "email_verification",
		UserType: "customer",
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("id = %d, want 0", got)
	}
}

func TestActivationCodeRepositoryFindByEmailAndUserType(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 4, 10, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.June, 4, 9, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(activationCodeQueryFindByEmailAndUT).
		WithArgs("customer@example.com", "customer").
		WillReturnRows(sqlmockRows(activationCodeColumns).
			AddRow(activationCodeRow(7, "123456", "customer@example.com", "email_verification", "customer", expiresAt, nil, createdAt)...))

	got, err := NewActivationCodeRepository(db).FindByEmailAndUserType(context.Background(), "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("FindByEmailAndUserType returned error: %v", err)
	}
	assertActivationCode(t, got, 7, "123456", "customer@example.com", "email_verification", "customer", expiresAt, nil, createdAt)
}

func TestActivationCodeRepositoryFindByCodeReturnsUsedAtWhenPresent(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 4, 10, 0, 0, 0, time.UTC)
	usedAt := time.Date(2026, time.June, 4, 9, 30, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.June, 4, 9, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(activationCodeQueryFindByCode).
		WithArgs("123456").
		WillReturnRows(sqlmockRows(activationCodeColumns).
			AddRow(activationCodeRow(7, "123456", "customer@example.com", "email_verification", "customer", expiresAt, &usedAt, createdAt)...))

	got, err := NewActivationCodeRepository(db).FindByCode(context.Background(), "123456")
	if err != nil {
		t.Fatalf("FindByCode returned error: %v", err)
	}
	assertActivationCode(t, got, 7, "123456", "customer@example.com", "email_verification", "customer", expiresAt, &usedAt, createdAt)
}

func TestActivationCodeRepositoryFindByCodeReturnsNilWhenNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery(activationCodeQueryFindByCode).
		WithArgs("missing").
		WillReturnRows(sqlmockRows(activationCodeColumns))

	got, err := NewActivationCodeRepository(db).FindByCode(context.Background(), "missing")
	if err != nil {
		t.Fatalf("FindByCode returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("activation code = %#v, want nil", got)
	}
}

func TestActivationCodeRepositoryMarkAsUsedWithTx(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(activationCodeQueryMarkAsUsed).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	if err := NewActivationCodeRepository(db).MarkAsUsedWithTx(context.Background(), tx, 7); err != nil {
		t.Fatalf("MarkAsUsedWithTx returned error: %v", err)
	}
}

func TestActivationCodeRepositoryMarkAsUsedWithTxReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("update failed")
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(activationCodeQueryMarkAsUsed).
		WithArgs(int64(7)).
		WillReturnError(wantErr)
	mock.ExpectRollback()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback()

	err = NewActivationCodeRepository(db).MarkAsUsedWithTx(context.Background(), tx, 7)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
}

func TestActivationCodeRepositoryDeleteExpiredCodes(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(activationCodeQueryDeleteExpired).
		WillReturnResult(sqlmock.NewResult(0, 3))

	if err := NewActivationCodeRepository(db).DeleteExpiredCodes(context.Background()); err != nil {
		t.Fatalf("DeleteExpiredCodes returned error: %v", err)
	}
}

func TestActivationCodeRepositoryDeleteExpiredCodesReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("delete failed")
	db, mock := newMockDB(t)
	mock.ExpectExec(activationCodeQueryDeleteExpired).WillReturnError(wantErr)

	err := NewActivationCodeRepository(db).DeleteExpiredCodes(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
}

func activationCodeRow(id int64, code, email, codeType, userType string, expiresAt time.Time, usedAt *time.Time, createdAt time.Time) []driver.Value {
	var used any
	if usedAt != nil {
		used = *usedAt
	}
	return []driver.Value{id, code, email, codeType, userType, expiresAt, used, createdAt}
}

func assertActivationCode(t *testing.T, got *model.ActivationCode, id int64, code, email, codeType, userType string, expiresAt time.Time, usedAt *time.Time, createdAt time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("activation code is nil")
	}
	if got.ID != id || got.Code != code || got.Email != email || got.Type != codeType ||
		got.UserType != userType || !got.ExpiresAt.Equal(expiresAt) || !got.CreatedAt.Equal(createdAt) {
		t.Fatalf("activation code = %#v", got)
	}
	assertTimePtr(t, "used at", got.UsedAt, usedAt)
}

func assertTimePtr(t *testing.T, field string, got *time.Time, want *time.Time) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("%s = %v, want nil", field, *got)
		}
		return
	}
	if got == nil || !got.Equal(*want) {
		t.Fatalf("%s = %v, want %v", field, got, *want)
	}
}
