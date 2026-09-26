package middleware

import (
	"net/http"
	"strings"

	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		scheme, token, ok := strings.Cut(ctx.GetHeader("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized("missing bearer token"))
			return
		}

		claims, err := util.ParseAccessToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized("invalid bearer token"))
			return
		}

		ctx.Set(UserIDKey, claims.UserID)
		ctx.Next()
	}
}
