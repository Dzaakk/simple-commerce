package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"Dzaakk/simple-commerce/package/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateActivationCode(t *testing.T) {
	code, err := generateActivationCode()
	if err != nil {
		t.Fatalf("generateActivationCode returned error: %v", err)
	}
	if len(code) != activationCodeBytes*2 {
		t.Fatalf("code length = %d, want %d", len(code), activationCodeBytes*2)
	}
	for _, r := range code {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			t.Fatalf("code contains non-hex character %q", r)
		}
	}
}

func TestHashPassword(t *testing.T) {
	got, err := hashPassword("  secret-password  ")
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}
	if got == "secret-password" {
		t.Fatal("hashPassword returned the plain password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(got), []byte("secret-password")); err != nil {
		t.Fatalf("hash does not match trimmed password: %v", err)
	}
}

func TestHashPasswordRejectsBlankPassword(t *testing.T) {
	_, err := hashPassword("   ")

	var appErr *response.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want *response.AppError", err, err)
	}
	if appErr.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want %d", appErr.Code, http.StatusBadRequest)
	}
	if appErr.Message != "password is required" {
		t.Fatalf("message = %q, want %q", appErr.Message, "password is required")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	tokenString, err := generateAccessToken("user-1", "customer", "customer@example.com")
	if err != nil {
		t.Fatalf("generateAccessToken returned error: %v", err)
	}

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !token.Valid {
		t.Fatal("token is not valid")
	}
	if claims.UserID != "user-1" || claims.UserType != "customer" || claims.Email != "customer@example.com" {
		t.Fatalf("claims = %#v", claims)
	}
	if claims.ExpiresAt == nil || time.Until(claims.ExpiresAt.Time) <= 14*time.Minute {
		t.Fatalf("expires_at = %v, want about %v from now", claims.ExpiresAt, accessTokenDuration)
	}
}

func TestGenerateAccessTokenRequiresSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := generateAccessToken("user-1", "customer", "customer@example.com")
	if err == nil {
		t.Fatal("generateAccessToken returned nil error")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	raw, hashed, err := generateRefreshToken()
	if err != nil {
		t.Fatalf("generateRefreshToken returned error: %v", err)
	}
	if len(raw) != refreshTokenBytes*2 {
		t.Fatalf("raw token length = %d, want %d", len(raw), refreshTokenBytes*2)
	}
	if len(hashed) != 64 {
		t.Fatalf("hashed token length = %d, want 64", len(hashed))
	}
	if hashRefreshToken(raw) != hashed {
		t.Fatalf("hashRefreshToken(raw) = %q, want %q", hashRefreshToken(raw), hashed)
	}
}
