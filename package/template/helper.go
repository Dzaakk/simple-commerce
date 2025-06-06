package template

import (
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizedChecker(ctx *gin.Context, customerId string) {
	currentId, exists := ctx.Get("customerId")
	if !exists || currentId != customerId {
		ctx.JSON(http.StatusUnauthorized, response.Unauthorized(nil))
		ctx.Abort()
		return
	}
}
