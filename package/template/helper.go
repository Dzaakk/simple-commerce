package template

import (
	"github.com/gin-gonic/gin"
)

func AuthorizedChecker(ctx *gin.Context, customerID string) bool {
	currentID, exists := ctx.Get("customerId")
	if !exists || currentID != customerID {
		return false
	}
	return true
}
