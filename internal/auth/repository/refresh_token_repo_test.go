package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/package/constant"
	"github.com/DATA-DOG/go-sqlmock"
)

var refreshTokenColumns = []string{"id", "user_id", "user_type", "token_hash", "expires_at", "revoked_at", "created_at"}

func TestRefreshTokenRepositoryCreate(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 4, 10, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.June, 4, 9, 0, 0, 0, time.UTC)
	token := &model.RefreshToken{
		UserID:    "customer-1",
		UserType:  constant.Customer,
		TokenHash: "token-hash",
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}
	db, mock := newMockDB(t)
	mock.ExpectQuery(refreshTokenQueryCreate).
		WithArgs(token.UserID, string(token.UserType), token.TokenHash, token.ExpiresAt, nil, token.CreatedAt).
		WillReturnRows(sqlmockRows([]string{"id"}).AddRow(int64(11)))

	got, err := NewRefreshTokenRepository(db).Create(context.Background(), token)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != 11 {
		t.Fatalf("id = %d, want 11", got)
	}
}

func TestRefreshTokenRepositoryCreateReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("insert failed")
	db, mock := newMockDB(t)
	mock.ExpectQuery(refreshTokenQueryCreate).
		WithArgs("customer-1", string(constant.Customer), "token-hash", time.Time{}, nil, time.Time{}).
		WillReturnError(wantErr)

	got, err := NewRefreshTokenRepository(db).Create(context.Background(), &model.RefreshToken{
		UserID:    "customer-1",
		UserType:  constant.Customer,
		TokenHash: "token-hash",
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("id = %d, want 0", got)
	}
}

func TestRefreshTokenRepositoryFindByUserID(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 4, 10, 0, 0, 0, time.UTC)
	revokedAt := time.Date(2026, time.June, 4, 9, 30, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.June, 4, 9, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(refreshTokenQueryFindByUser).
		WithArgs("customer-1").
		WillReturnRows(sqlmockRows(refreshTokenColumns).
			AddRow(refreshTokenRow(11, "customer-1", constant.Customer, "token-hash", expiresAt, &revokedAt, createdAt)...))

	got, err := NewRefreshTokenRepository(db).FindByUserID(context.Background(), "customer-1")
	if err != nil {
		t.Fatalf("FindByUserID returned error: %v", err)
	}
	assertRefreshToken(t, got, 11, "customer-1", constant.Customer, "token-hash", expiresAt, &revokedAt, createdAt)
}

func TestRefreshTokenRepositoryFindByTokenHashReturnsNilWhenNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery(refreshTokenQueryFindByTokenHash).
		WithArgs("missing-hash").
		WillReturnRows(sqlmockRows(refreshTokenColumns))

	got, err := NewRefreshTokenRepository(db).FindByTokenHash(context.Background(), "missing-hash")
	if err != nil {
		t.Fatalf("FindByTokenHash returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("refresh token = %#v, want nil", got)
	}
}

func TestRefreshTokenRepositorySetExpire(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 5, 10, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQuerySetExpire).
		WithArgs(expiresAt, int64(11)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	got, err := NewRefreshTokenRepository(db).SetExpire(context.Background(), 11, expiresAt)
	if err != nil {
		t.Fatalf("SetExpire returned error: %v", err)
	}
	if got != 1 {
		t.Fatalf("rows affected = %d, want 1", got)
	}
}

func TestRefreshTokenRepositorySetExpireReturnsNoRowsError(t *testing.T) {
	expiresAt := time.Date(2026, time.June, 5, 10, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQuerySetExpire).
		WithArgs(expiresAt, int64(11)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	got, err := NewRefreshTokenRepository(db).SetExpire(context.Background(), 11, expiresAt)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
	if got != 0 {
		t.Fatalf("rows affected = %d, want 0", got)
	}
}

func TestRefreshTokenRepositoryRevoke(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQueryRevoke).
		WithArgs("token-hash").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := NewRefreshTokenRepository(db).Revoke(context.Background(), "token-hash"); err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}
}

func TestRefreshTokenRepositoryRevokeReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("update failed")
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQueryRevoke).
		WithArgs("token-hash").
		WillReturnError(wantErr)

	err := NewRefreshTokenRepository(db).Revoke(context.Background(), "token-hash")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
}

func TestRefreshTokenRepositoryRevokeAllByUser(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQueryRevokeAllByUser).
		WithArgs("customer-1", string(constant.Customer)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	if err := NewRefreshTokenRepository(db).RevokeAllByUser(context.Background(), "customer-1", constant.Customer); err != nil {
		t.Fatalf("RevokeAllByUser returned error: %v", err)
	}
}

func TestRefreshTokenRepositoryRevokeAllByUserReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("update failed")
	db, mock := newMockDB(t)
	mock.ExpectExec(refreshTokenQueryRevokeAllByUser).
		WithArgs("customer-1", string(constant.Customer)).
		WillReturnError(wantErr)

	err := NewRefreshTokenRepository(db).RevokeAllByUser(context.Background(), "customer-1", constant.Customer)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
}

func refreshTokenRow(id int64, userID string, userType constant.UserType, tokenHash string, expiresAt time.Time, revokedAt *time.Time, createdAt time.Time) []driver.Value {
	var revoked any
	if revokedAt != nil {
		revoked = *revokedAt
	}
	return []driver.Value{id, userID, string(userType), tokenHash, expiresAt, revoked, createdAt}
}

func assertRefreshToken(t *testing.T, got *model.RefreshToken, id int64, userID string, userType constant.UserType, tokenHash string, expiresAt time.Time, revokedAt *time.Time, createdAt time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("refresh token is nil")
	}
	if got.ID != id || got.UserID != userID || got.UserType != userType || got.TokenHash != tokenHash ||
		!got.ExpiresAt.Equal(expiresAt) || !got.CreatedAt.Equal(createdAt) {
		t.Fatalf("refresh token = %#v", got)
	}
	assertTimePtr(t, "revoked at", got.RevokedAt, revokedAt)
}
