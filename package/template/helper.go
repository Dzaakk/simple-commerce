package template

import (
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
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response(http.StatusUnauthorized, "Unauthorized", "User not authorized"))
		ctx.Abort()
		return
	}
	if currentId != customerId {
		ctx.JSON(http.StatusInternalServerError, Response(http.StatusInternalServerError, "Internal Server Error", "Internal Server Error"))
		ctx.Abort()
		return
	}
}
