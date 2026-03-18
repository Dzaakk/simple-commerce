package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	activationCodeBytes  = 16 // 32 hex chars
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
	refreshTokenBytes    = 32
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func generateActivationCode() (string, error) {
	b := make([]byte, activationCodeBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate activation code: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashPassword(plain string) (string, error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return "", errors.New("password is required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hashed), nil
}

func generateAccessToken(userID, userType, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	claims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

func generateRefreshToken() (raw string, hashed string, err error) {
	b := make([]byte, refreshTokenBytes)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hashed = hex.EncodeToString(sum[:])

	return raw, hashed, nil
}

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
