package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	registerCustomerFn func(context.Context, *dto.RegisterCustomerRequest) error
	registerSellerFn   func(context.Context, *dto.RegisterSellerRequest) error
	verifyEmailFn      func(context.Context, string) error
	loginFn            func(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	refreshTokenFn     func(context.Context, string) (*dto.RefreshTokenResponse, error)
	logoutFn           func(context.Context, string) error
}

func (m *mockAuthService) RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerRequest) error {
	if m.registerCustomerFn == nil {
		return errors.New("unexpected RegisterCustomer call")
	}
	return m.registerCustomerFn(ctx, req)
}

func (m *mockAuthService) RegisterSeller(ctx context.Context, req *dto.RegisterSellerRequest) error {
	if m.registerSellerFn == nil {
		return errors.New("unexpected RegisterSeller call")
	}
	return m.registerSellerFn(ctx, req)
}

func (m *mockAuthService) VerifyEmail(ctx context.Context, activationCode string) error {
	if m.verifyEmailFn == nil {
		return errors.New("unexpected VerifyEmail call")
	}
	return m.verifyEmailFn(ctx, activationCode)
}

func (m *mockAuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if m.loginFn == nil {
		return nil, errors.New("unexpected Login call")
	}
	return m.loginFn(ctx, req)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, rawRefreshToken string) (*dto.RefreshTokenResponse, error) {
	if m.refreshTokenFn == nil {
		return nil, errors.New("unexpected RefreshToken call")
	}
	return m.refreshTokenFn(ctx, rawRefreshToken)
}

func (m *mockAuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	if m.logoutFn == nil {
		return errors.New("unexpected Logout call")
	}
	return m.logoutFn(ctx, rawRefreshToken)
}

func TestAuthHandlerLogin(t *testing.T) {
	handler := NewAuthHandler(&mockAuthService{
		loginFn: func(_ context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
			if req.Email != "customer@example.com" || req.Password != "secret" || req.UserType != "customer" {
				t.Fatalf("login request = %#v", req)
			}
			return &dto.LoginResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    900,
			}, nil
		},
	})

	w, ctx := performAuthRequest(http.MethodPost, "/login", `{"email":"customer@example.com","password":"secret","user_type":"customer"}`, handler.Login)
	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got response.ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	data, ok := got.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data = %T %#v, want map", got.Data, got.Data)
	}
	if data["AccessToken"] != "access-token" || data["RefreshToken"] != "refresh-token" {
		t.Fatalf("data = %#v", data)
	}
}

func TestAuthHandlerLoginRejectsInvalidRequest(t *testing.T) {
	called := false
	handler := NewAuthHandler(&mockAuthService{
		loginFn: func(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error) {
			called = true
			return nil, nil
		},
	})

	_, ctx := performAuthRequest(http.MethodPost, "/login", `{"email":"not-an-email","password":"secret","user_type":"customer"}`, handler.Login)
	if called {
		t.Fatal("service Login must not be called for invalid request")
	}
	assertHandlerAppError(t, ctx, http.StatusBadRequest, "invalid request data")
}

func TestAuthHandlerVerifyEmail(t *testing.T) {
	handler := NewAuthHandler(&mockAuthService{
		verifyEmailFn: func(_ context.Context, activationCode string) error {
			if activationCode != "activation-code" {
				t.Fatalf("activation code = %q, want activation-code", activationCode)
			}
			return nil
		},
	})

	w, ctx := performAuthRequest(http.MethodGet, "/verify-email?code=activation-code", "", handler.VerifyEmail)
	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuthHandlerVerifyEmailRejectsMissingCode(t *testing.T) {
	called := false
	handler := NewAuthHandler(&mockAuthService{
		verifyEmailFn: func(context.Context, string) error {
			called = true
			return nil
		},
	})

	_, ctx := performAuthRequest(http.MethodGet, "/verify-email", "", handler.VerifyEmail)
	if called {
		t.Fatal("service VerifyEmail must not be called without code")
	}
	assertHandlerAppError(t, ctx, http.StatusBadRequest, "invalid request data")
}

func TestAuthHandlerRefreshToken(t *testing.T) {
	handler := NewAuthHandler(&mockAuthService{
		refreshTokenFn: func(_ context.Context, rawRefreshToken string) (*dto.RefreshTokenResponse, error) {
			if rawRefreshToken != "refresh-token" {
				t.Fatalf("refresh token = %q, want refresh-token", rawRefreshToken)
			}
			return &dto.RefreshTokenResponse{AccessToken: "new-access-token", ExpiresIn: 900}, nil
		},
	})

	w, ctx := performAuthRequest(http.MethodPost, "/refresh-token", `{"refresh_token":"refresh-token"}`, handler.RefreshToken)
	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "new-access-token") {
		t.Fatalf("body = %q, want access token", w.Body.String())
	}
}

func TestAuthHandlerLogout(t *testing.T) {
	handler := NewAuthHandler(&mockAuthService{
		logoutFn: func(_ context.Context, rawRefreshToken string) error {
			if rawRefreshToken != "refresh-token" {
				t.Fatalf("refresh token = %q, want refresh-token", rawRefreshToken)
			}
			return nil
		},
	})

	w, ctx := performAuthRequest(http.MethodPost, "/logout", `{"refresh_token":"refresh-token"}`, handler.Logout)
	if len(ctx.Errors) != 0 {
		t.Fatalf("errors = %v, want none", ctx.Errors)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func performAuthRequest(method, target, body string, handlerFunc gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	handlerFunc(ctx)

	return w, ctx
}

func assertHandlerAppError(t *testing.T, ctx *gin.Context, code int, message string) {
	t.Helper()

	if len(ctx.Errors) != 1 {
		t.Fatalf("errors = %v, want exactly one", ctx.Errors)
	}
	var appErr *response.AppError
	if !errors.As(ctx.Errors[0].Err, &appErr) {
		t.Fatalf("error = %T %v, want *response.AppError", ctx.Errors[0].Err, ctx.Errors[0].Err)
	}
	if appErr.Code != code {
		t.Fatalf("code = %d, want %d", appErr.Code, code)
	}
	if appErr.Message != message {
		t.Fatalf("message = %q, want %q", appErr.Message, message)
	}
}
