package util

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratePlaceHolders(n int) string {
	holders := make([]string, n)

	for i := 1; i <= n; i++ {
		holders[i-1] = fmt.Sprintf("$%d", i)
	}

	return strings.Join(holders, ", ")
}

func AuthorizedChecker(ctx *gin.Context, customerID string) bool {
	currentID, exists := ctx.Get("customerId")
	if !exists || currentID != customerID {
		return false
	}
	return true
}
