package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const activationCodeBytes = 16 // 32 hex chars

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
