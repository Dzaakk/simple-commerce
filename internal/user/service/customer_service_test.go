package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
)

type mockCustomerRepository struct {
	createFn       func(context.Context, *model.Customer) (string, error)
	updateFn       func(context.Context, *model.Customer) (int64, error)
	findByIDFn     func(context.Context, string) (*model.Customer, error)
	findByEmailFn  func(context.Context, string) (*model.Customer, error)
	updateStatusFn func(context.Context, string, constant.UserStatus) error
}

func (f *mockCustomerRepository) Create(ctx context.Context, data *model.Customer) (string, error) {
	if f.createFn == nil {
		return "", errors.New("unexpected Create call")
	}
	return f.createFn(ctx, data)
}

func (f *mockCustomerRepository) Update(ctx context.Context, data *model.Customer) (int64, error) {
	if f.updateFn == nil {
		return 0, errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, data)
}

func (f *mockCustomerRepository) FindByID(ctx context.Context, customerID string) (*model.Customer, error) {
	if f.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return f.findByIDFn(ctx, customerID)
}

func (f *mockCustomerRepository) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	if f.findByEmailFn == nil {
		return nil, errors.New("unexpected FindByEmail call")
	}
	return f.findByEmailFn(ctx, email)
}

func (f *mockCustomerRepository) UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error {
	if f.updateStatusFn == nil {
		return errors.New("unexpected UpdateStatus call")
	}
	return f.updateStatusFn(ctx, customerID, status)
}

func TestCustomerServiceCreate(t *testing.T) {
	ctx := context.Background()
	req := &dto.RegisterCustomerRequest{
		Email:    "customer@example.com",
		Password: "plain-password",
		FullName: "Customer Name",
		Phone:    "08123456789",
	}

	repo := &mockCustomerRepository{
		createFn: func(_ context.Context, data *model.Customer) (string, error) {
			if data.Email != req.Email {
				t.Fatalf("email = %q, want %q", data.Email, req.Email)
			}
			if data.PasswordHash != req.Password {
				t.Fatalf("password hash = %q, want %q", data.PasswordHash, req.Password)
			}
			if data.FullName != req.FullName {
				t.Fatalf("full name = %q, want %q", data.FullName, req.FullName)
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
			return "customer-1", nil
		},
	}

	got, err := NewCustomerService(repo).Create(ctx, req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != "customer-1" {
		t.Fatalf("id = %q, want %q", got, "customer-1")
	}
}

func TestCustomerServiceCreateReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("database down")
	repo := &mockCustomerRepository{
		createFn: func(context.Context, *model.Customer) (string, error) {
			return "", wantErr
		},
	}

	got, err := NewCustomerService(repo).Create(context.Background(), &dto.RegisterCustomerRequest{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if got != "" {
		t.Fatalf("id = %q, want empty", got)
	}
}

func TestCustomerServiceUpdateRejectsInvalidCustomerID(t *testing.T) {
	tests := []struct {
		name       string
		customerID string
	}{
		{name: "not numeric", customerID: "abc"},
		{name: "zero", customerID: "0"},
		{name: "negative", customerID: "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			repo := &mockCustomerRepository{
				updateFn: func(context.Context, *model.Customer) (int64, error) {
					called = true
					return 1, nil
				},
			}

			err := NewCustomerService(repo).Update(context.Background(), &dto.UpdateReq{CustomerID: tt.customerID})
			assertAppError(t, err, http.StatusBadRequest, "invalid parameter customer id")
			if called {
				t.Fatal("repository Update must not be called for invalid customer id")
			}
		})
	}
}

func TestCustomerServiceUpdate(t *testing.T) {
	req := &dto.UpdateReq{
		CustomerID: "42",
		Email:      "new@example.com",
		FullName:   "New Name",
		Phone:      "08999999999",
	}
	repo := &mockCustomerRepository{
		updateFn: func(_ context.Context, data *model.Customer) (int64, error) {
			if data.ID != "42" {
				t.Fatalf("id = %q, want %q", data.ID, "42")
			}
			if data.Email != req.Email {
				t.Fatalf("email = %q, want %q", data.Email, req.Email)
			}
			if data.FullName != req.FullName {
				t.Fatalf("full name = %q, want %q", data.FullName, req.FullName)
			}
			if data.Phone != req.Phone {
				t.Fatalf("phone = %q, want %q", data.Phone, req.Phone)
			}
			if data.Status != string(constant.StatusPending) {
				t.Fatalf("status = %q, want default %q", data.Status, constant.StatusPending)
			}
			if data.UpdatedAt.IsZero() {
				t.Fatal("updated_at must be set")
			}
			return 1, nil
		},
	}

	if err := NewCustomerService(repo).Update(context.Background(), req); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
}

func TestCustomerServiceUpdateReturnsNotFoundWhenNoRowsAffected(t *testing.T) {
	repo := &mockCustomerRepository{
		updateFn: func(context.Context, *model.Customer) (int64, error) {
			return 0, nil
		},
	}

	err := NewCustomerService(repo).Update(context.Background(), &dto.UpdateReq{CustomerID: "42"})
	assertAppError(t, err, http.StatusNotFound, "customer not found")
}

func TestCustomerServiceUpdateReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("write failed")
	repo := &mockCustomerRepository{
		updateFn: func(context.Context, *model.Customer) (int64, error) {
			return 0, wantErr
		},
	}

	err := NewCustomerService(repo).Update(context.Background(), &dto.UpdateReq{CustomerID: "42"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestCustomerServiceFindByEmail(t *testing.T) {
	want := &model.Customer{ID: "customer-1", Email: "customer@example.com"}
	repo := &mockCustomerRepository{
		findByEmailFn: func(_ context.Context, email string) (*model.Customer, error) {
			if email != want.Email {
				t.Fatalf("email = %q, want %q", email, want.Email)
			}
			return want, nil
		},
	}

	got, err := NewCustomerService(repo).FindByEmail(context.Background(), want.Email)
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	if got != want {
		t.Fatalf("customer = %#v, want same pointer %#v", got, want)
	}
}

func TestCustomerServiceFindByIDMapsResponse(t *testing.T) {
	now := time.Date(2026, time.May, 24, 10, 0, 0, 0, time.UTC)
	repo := &mockCustomerRepository{
		findByIDFn: func(_ context.Context, customerID string) (*model.Customer, error) {
			if customerID != "customer-1" {
				t.Fatalf("customer id = %q, want %q", customerID, "customer-1")
			}
			return &model.Customer{
				ID:        "customer-1",
				Email:     "customer@example.com",
				FullName:  "Customer Name",
				Phone:     "08123456789",
				Status:    string(constant.StatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	got, err := NewCustomerService(repo).FindByID(context.Background(), "customer-1")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if got.ID != "customer-1" || got.Email != "customer@example.com" || got.FullName != "Customer Name" ||
		got.Phone != "08123456789" || got.Status != string(constant.StatusActive) ||
		!got.CreatedAt.Equal(now) || !got.UpdatedAt.Equal(now) {
		t.Fatalf("customer response = %#v", got)
	}
}

func TestCustomerServiceFindByIDReturnsNilWhenNotFound(t *testing.T) {
	repo := &mockCustomerRepository{
		findByIDFn: func(context.Context, string) (*model.Customer, error) {
			return nil, nil
		},
	}

	got, err := NewCustomerService(repo).FindByID(context.Background(), "missing")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("customer = %#v, want nil", got)
	}
}

func TestCustomerServiceUpdateStatusDelegatesToRepository(t *testing.T) {
	repo := &mockCustomerRepository{
		updateStatusFn: func(_ context.Context, customerID string, status constant.UserStatus) error {
			if customerID != "customer-1" {
				t.Fatalf("customer id = %q, want %q", customerID, "customer-1")
			}
			if status != constant.StatusActive {
				t.Fatalf("status = %q, want %q", status, constant.StatusActive)
			}
			return nil
		},
	}

	if err := NewCustomerService(repo).UpdateStatus(context.Background(), "customer-1", constant.StatusActive); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
}

func TestCustomerServiceUpdateStatusReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("transaction failed")
	repo := &mockCustomerRepository{
		updateStatusFn: func(context.Context, string, constant.UserStatus) error {
			return wantErr
		},
	}

	err := NewCustomerService(repo).UpdateStatus(context.Background(), "customer-1", constant.StatusActive)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func assertAppError(t *testing.T, err error, code int, message string) {
	t.Helper()

	var appErr *response.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want *response.AppError", err, err)
	}
	if appErr.Code != code {
		t.Fatalf("code = %d, want %d", appErr.Code, code)
	}
	if appErr.Message != message {
		t.Fatalf("message = %q, want %q", appErr.Message, message)
	}
}
