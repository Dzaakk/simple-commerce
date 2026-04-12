package middleware

import (
	"Dzaakk/simple-commerce/package/response"
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler centralizes error responses.
// Handlers should call ctx.Error(err) and return.
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}

		lastErr := ctx.Errors.Last()
		if lastErr == nil {
			return
		}

		status, res := response.ErrorResponse(lastErr.Err)
		if status >= 500 {
			log.Printf("error: %s %s: %v", ctx.Request.Method, ctx.Request.URL.Path, lastErr.Err)
		}

		ctx.JSON(status, res)
	}
}
