package template

import (
	"Dzaakk/simple-commerce/package/response"
	"net/http"

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

func AuthorizedChecker(ctx *gin.Context, customerId string) {
	currentId, exists := ctx.Get("customerId")
	if !exists || currentId != customerId {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(nil))
		ctx.Abort()
		return
	}
}
