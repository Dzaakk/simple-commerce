package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
)

type mockSellerRepository struct {
	createFn             func(context.Context, *model.Seller) (string, error)
	updateFn             func(context.Context, *model.Seller) (int64, error)
	findByIDFn           func(context.Context, string) (*model.Seller, error)
	findByEmailFn        func(context.Context, string) (*model.Seller, error)
	findByShopNameFn     func(context.Context, string) ([]*model.Seller, error)
	updateStatusFn       func(context.Context, string, constant.UserStatus) error
	updateStatusWithTxFn func(context.Context, *sql.Tx, string, constant.UserStatus) error
}

func (f *mockSellerRepository) Create(ctx context.Context, data *model.Seller) (string, error) {
	if f.createFn == nil {
		return "", errors.New("unexpected Create call")
	}
	return f.createFn(ctx, data)
}

func (f *mockSellerRepository) Update(ctx context.Context, data *model.Seller) (int64, error) {
	if f.updateFn == nil {
		return 0, errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, data)
}

func (f *mockSellerRepository) FindByID(ctx context.Context, sellerID string) (*model.Seller, error) {
	if f.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return f.findByIDFn(ctx, sellerID)
}

func (f *mockSellerRepository) FindByEmail(ctx context.Context, email string) (*model.Seller, error) {
	if f.findByEmailFn == nil {
		return nil, errors.New("unexpected FindByEmail call")
	}
	return f.findByEmailFn(ctx, email)
}

func (f *mockSellerRepository) FindByShopName(ctx context.Context, name string) ([]*model.Seller, error) {
	if f.findByShopNameFn == nil {
		return nil, errors.New("unexpected FindByShopName call")
	}
	return f.findByShopNameFn(ctx, name)
}

func (f *mockSellerRepository) UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error {
	if f.updateStatusFn == nil {
		return errors.New("unexpected UpdateStatus call")
	}
	return f.updateStatusFn(ctx, sellerID, status)
}

func (f *mockSellerRepository) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error {
	if f.updateStatusWithTxFn == nil {
		return errors.New("unexpected UpdateStatusWithTx call")
	}
	return f.updateStatusWithTxFn(ctx, tx, sellerID, status)
}

func TestSellerServiceCreate(t *testing.T) {
	req := &dto.RegisterSellerRequest{
		Email:    "seller@example.com",
		Password: "plain-password",
		FullName: "Seller Name",
		Phone:    "08123456789",
		ShopName: "Good Shop",
	}
	repo := &mockSellerRepository{
		createFn: func(_ context.Context, data *model.Seller) (string, error) {
			if data.Email != req.Email {
				t.Fatalf("email = %q, want %q", data.Email, req.Email)
			}
			if data.PasswordHash != req.Password {
				t.Fatalf("password hash = %q, want %q", data.PasswordHash, req.Password)
			}
			if data.ShopName != req.ShopName {
				t.Fatalf("shop name = %q, want %q", data.ShopName, req.ShopName)
			}
			if data.Phone != req.Phone {
				t.Fatalf("phone = %q, want %q", data.Phone, req.Phone)
			}
			if data.Status != string(constant.StatusPending) {
				t.Fatalf("status = %q, want %q", data.Status, constant.StatusPending)
			}
			if data.CreatedAt.IsZero() || data.UpdatedAt.IsZero() {
				t.Fatal("created_at and updated_at must be set")
			}
			return "seller-1", nil
		},
	}

	got, err := NewSellerService(repo).Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != "seller-1" {
		t.Fatalf("id = %q, want %q", got, "seller-1")
	}
}

func TestSellerServiceCreateReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("database down")
	repo := &mockSellerRepository{
		createFn: func(context.Context, *model.Seller) (string, error) {
			return "", wantErr
		},
	}

	got, err := NewSellerService(repo).Create(context.Background(), &dto.RegisterSellerRequest{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if got != "" {
		t.Fatalf("id = %q, want empty", got)
	}
}

func TestSellerServiceUpdateRejectsInvalidSellerID(t *testing.T) {
	tests := []struct {
		name     string
		sellerID string
	}{
		{name: "not numeric", sellerID: "abc"},
		{name: "zero", sellerID: "0"},
		{name: "negative", sellerID: "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			repo := &mockSellerRepository{
				updateFn: func(context.Context, *model.Seller) (int64, error) {
					called = true
					return 1, nil
				},
			}

			err := NewSellerService(repo).Update(context.Background(), &dto.SellerUpdateReq{SellerID: tt.sellerID})
			assertAppError(t, err, http.StatusBadRequest, "invalid parameter seller id")
			if called {
				t.Fatal("repository Update must not be called for invalid seller id")
			}
		})
	}
}

func TestSellerServiceUpdate(t *testing.T) {
	req := &dto.SellerUpdateReq{
		SellerID: "42",
		Email:    "new-seller@example.com",
		ShopName: "Better Shop",
		Phone:    "08999999999",
		Status:   string(constant.StatusActive),
	}
	repo := &mockSellerRepository{
		updateFn: func(_ context.Context, data *model.Seller) (int64, error) {
			if data.ID != "42" {
				t.Fatalf("id = %q, want %q", data.ID, "42")
			}
			if data.Email != req.Email {
				t.Fatalf("email = %q, want %q", data.Email, req.Email)
			}
			if data.ShopName != req.ShopName {
				t.Fatalf("shop name = %q, want %q", data.ShopName, req.ShopName)
			}
			if data.Phone != req.Phone {
				t.Fatalf("phone = %q, want %q", data.Phone, req.Phone)
			}
			if data.Status != string(constant.StatusActive) {
				t.Fatalf("status = %q, want %q", data.Status, constant.StatusActive)
			}
			if data.UpdatedAt.IsZero() {
				t.Fatal("updated_at must be set")
			}
			return 1, nil
		},
	}

	if err := NewSellerService(repo).Update(context.Background(), req); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
}

func TestSellerServiceUpdateReturnsNotFoundWhenNoRowsAffected(t *testing.T) {
	repo := &mockSellerRepository{
		updateFn: func(context.Context, *model.Seller) (int64, error) {
			return 0, nil
		},
	}

	err := NewSellerService(repo).Update(context.Background(), &dto.SellerUpdateReq{SellerID: "42"})
	assertAppError(t, err, http.StatusNotFound, "seller not found")
}

func TestSellerServiceFindByEmail(t *testing.T) {
	want := &model.Seller{ID: "seller-1", Email: "seller@example.com"}
	repo := &mockSellerRepository{
		findByEmailFn: func(_ context.Context, email string) (*model.Seller, error) {
			if email != want.Email {
				t.Fatalf("email = %q, want %q", email, want.Email)
			}
			return want, nil
		},
	}

	got, err := NewSellerService(repo).FindByEmail(context.Background(), want.Email)
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	if got != want {
		t.Fatalf("seller = %#v, want same pointer %#v", got, want)
	}
}

func TestSellerServiceFindByIDMapsResponse(t *testing.T) {
	now := time.Date(2026, time.May, 24, 10, 0, 0, 0, time.UTC)
	repo := &mockSellerRepository{
		findByIDFn: func(_ context.Context, sellerID string) (*model.Seller, error) {
			if sellerID != "seller-1" {
				t.Fatalf("seller id = %q, want %q", sellerID, "seller-1")
			}
			return &model.Seller{
				ID:        "seller-1",
				Email:     "seller@example.com",
				ShopName:  "Good Shop",
				Phone:     "08123456789",
				Status:    string(constant.StatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	got, err := NewSellerService(repo).FindByID(context.Background(), "seller-1")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if got.ID != "seller-1" || got.Email != "seller@example.com" || got.ShopName != "Good Shop" ||
		got.Phone != "08123456789" || got.Status != string(constant.StatusActive) ||
		!got.CreatedAt.Equal(now) || !got.UpdatedAt.Equal(now) {
		t.Fatalf("seller response = %#v", got)
	}
}

func TestSellerServiceFindByShopName(t *testing.T) {
	now := time.Date(2026, time.May, 24, 10, 0, 0, 0, time.UTC)
	repo := &mockSellerRepository{
		findByShopNameFn: func(_ context.Context, name string) ([]*model.Seller, error) {
			if name != "shop" {
				t.Fatalf("shop name = %q, want %q", name, "shop")
			}
			return []*model.Seller{
				{
					ID:        "seller-1",
					Email:     "seller@example.com",
					ShopName:  "Good Shop",
					Phone:     "08123456789",
					Status:    string(constant.StatusActive),
					CreatedAt: now,
					UpdatedAt: now,
				},
				nil,
			}, nil
		},
	}

	got, err := NewSellerService(repo).FindByShopName(context.Background(), "shop")
	if err != nil {
		t.Fatalf("FindByShopName returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("seller count = %d, want 1", len(got))
	}
	if got[0].ID != "seller-1" || got[0].ShopName != "Good Shop" {
		t.Fatalf("seller response = %#v", got[0])
	}
}

func TestSellerServiceFindByShopNameReturnsEmptySliceWhenNotFound(t *testing.T) {
	repo := &mockSellerRepository{
		findByShopNameFn: func(context.Context, string) ([]*model.Seller, error) {
			return nil, nil
		},
	}

	got, err := NewSellerService(repo).FindByShopName(context.Background(), "missing")
	if err != nil {
		t.Fatalf("FindByShopName returned error: %v", err)
	}
	if got == nil {
		t.Fatal("sellers slice must not be nil")
	}
	if len(got) != 0 {
		t.Fatalf("seller count = %d, want 0", len(got))
	}
}

func TestSellerServiceUpdateStatusDelegatesToRepository(t *testing.T) {
	repo := &mockSellerRepository{
		updateStatusFn: func(_ context.Context, sellerID string, status constant.UserStatus) error {
			if sellerID != "seller-1" {
				t.Fatalf("seller id = %q, want %q", sellerID, "seller-1")
			}
			if status != constant.StatusActive {
				t.Fatalf("status = %q, want %q", status, constant.StatusActive)
			}
			return nil
		},
	}

	if err := NewSellerService(repo).UpdateStatus(context.Background(), "seller-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
}

func TestSellerServiceUpdateStatusWithTxReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("transaction failed")
	repo := &mockSellerRepository{
		updateStatusWithTxFn: func(context.Context, *sql.Tx, string, constant.UserStatus) error {
			return wantErr
		},
	}

	err := NewSellerService(repo).UpdateStatusWithTx(context.Background(), nil, "seller-1", constant.StatusActive)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
