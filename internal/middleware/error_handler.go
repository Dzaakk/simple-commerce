package middleware

import (
	"log"

	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 || ctx.Writer.Written() {
			return
		}

		err := ctx.Errors.Last().Err
		status, body := response.ErrorResponse(err)
		if status >= 500 {
			log.Printf("request failed: method=%s path=%s error=%v", ctx.Request.Method, ctx.Request.URL.Path, err)
		}
		ctx.JSON(status, body)
	}
}
