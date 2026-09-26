package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"Dzaakk/simple-commerce/package/response"

	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenDuration = 7 * 24 * time.Hour
	refreshTokenBytes    = 32
)

func hashPassword(plain string) (string, error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return "", response.NewAppError(http.StatusBadRequest, "password is required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func generateRefreshToken() (string, string, error) {
	random := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(random); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	raw := hex.EncodeToString(random)
	return raw, hashRefreshToken(raw), nil
}

func hashRefreshToken(raw string) string {
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
