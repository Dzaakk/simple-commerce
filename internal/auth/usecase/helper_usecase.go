package usecase

import (
	"math/rand"
	"time"
)

func generateActivationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, codeLength)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
