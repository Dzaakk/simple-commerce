package service

import (
	"context"
	"errors"
	"testing"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	usermodel "Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"

	"golang.org/x/crypto/bcrypt"
)

type customerStub struct {
	create      func(*userdto.RegisterCustomerRequest) (string, error)
	findByEmail func(string) (*usermodel.Customer, error)
	findByID    func(string) (*userdto.CustomerResponse, error)
}

func (s customerStub) Create(_ context.Context, req *userdto.RegisterCustomerRequest) (string, error) {
	return s.create(req)
}
func (s customerStub) FindByEmail(_ context.Context, email string) (*usermodel.Customer, error) {
	return s.findByEmail(email)
}
func (s customerStub) FindByID(_ context.Context, id string) (*userdto.CustomerResponse, error) {
	return s.findByID(id)
}

type tokenStub struct {
	create func(*model.RefreshToken) error
	find   func(string) (*model.RefreshToken, error)
	revoke func(string) error
}

func (s tokenStub) Create(_ context.Context, token *model.RefreshToken) error { return s.create(token) }
func (s tokenStub) FindByTokenHash(_ context.Context, hash string) (*model.RefreshToken, error) {
	return s.find(hash)
}
func (s tokenStub) Revoke(_ context.Context, hash string) error { return s.revoke(hash) }

func TestRegisterCustomerNormalizesAndHashesInput(t *testing.T) {
	customers := customerStub{
		findByEmail: func(email string) (*usermodel.Customer, error) {
			if email != "customer@example.com" {
				t.Fatalf("normalized email = %q", email)
			}
			return nil, nil
		},
		create: func(req *userdto.RegisterCustomerRequest) (string, error) {
			if req.Email != "customer@example.com" || req.FullName != "Customer Name" {
				t.Fatalf("create request = %#v", req)
			}
			if bcrypt.CompareHashAndPassword([]byte(req.PasswordHash), []byte("secret-pass")) != nil {
				t.Fatal("password was not hashed")
			}
			return "customer-1", nil
		},
	}

	err := NewAuthService(customers, nil).RegisterCustomer(context.Background(), &dto.RegisterCustomerRequest{
		Email: " Customer@Example.COM ", Password: "secret-pass", FullName: " Customer Name ",
	})
	if err != nil {
		t.Fatalf("RegisterCustomer() error = %v", err)
	}
}

func TestRegisterCustomerRejectsDuplicate(t *testing.T) {
	customers := customerStub{
		findByEmail: func(string) (*usermodel.Customer, error) { return &usermodel.Customer{}, nil },
		create:      func(*userdto.RegisterCustomerRequest) (string, error) { t.Fatal("unexpected create"); return "", nil },
	}
	err := NewAuthService(customers, nil).RegisterCustomer(context.Background(), &dto.RegisterCustomerRequest{
		Email: "customer@example.com", Password: "secret-pass", FullName: "Customer",
	})
	if !errors.Is(err, response.ErrEmailAlreadyExist) {
		t.Fatalf("error = %v", err)
	}
}

func TestLoginCreatesUsableTokens(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-sufficient-entropy")
	hash, err := hashPassword("secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	customers := customerStub{
		findByEmail: func(string) (*usermodel.Customer, error) {
			return &usermodel.Customer{ID: "customer-1", Email: "customer@example.com", PasswordHash: hash, Status: string(constant.StatusActive)}, nil
		},
	}
	var stored *model.RefreshToken
	tokens := tokenStub{create: func(token *model.RefreshToken) error { stored = token; return nil }}

	result, err := NewAuthService(customers, tokens).Login(context.Background(), &dto.LoginRequest{Email: "customer@example.com", Password: "secret-pass"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	claims, err := util.ParseAccessToken(result.AccessToken)
	if err != nil || claims.UserID != "customer-1" {
		t.Fatalf("claims = %#v, error = %v", claims, err)
	}
	if stored == nil || stored.TokenHash != hashRefreshToken(result.RefreshToken) {
		t.Fatalf("stored refresh token = %#v", stored)
	}
}

func TestRefreshRejectsUnknownToken(t *testing.T) {
	tokens := tokenStub{find: func(string) (*model.RefreshToken, error) { return nil, nil }}
	_, err := NewAuthService(nil, tokens).RefreshToken(context.Background(), "missing")
	if !errors.Is(err, response.ErrInvalidRefreshToken) {
		t.Fatalf("error = %v", err)
	}
}

func TestLogoutHashesPresentedToken(t *testing.T) {
	tokens := tokenStub{revoke: func(hash string) error {
		if hash != hashRefreshToken("raw-token") {
			t.Fatalf("hash = %q", hash)
		}
		return nil
	}}
	if err := NewAuthService(nil, tokens).Logout(context.Background(), "raw-token"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
}
