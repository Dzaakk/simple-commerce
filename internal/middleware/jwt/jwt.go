package middleware

import (
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	AuthCache *repository.AuthCache
}

func NewJwtMiddleware(authCache *repository.AuthCache) *JWTMiddleware {
	return &JWTMiddleware{AuthCache: authCache}
}

func (m *JWTMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Authorization") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized("missing token"))
			return
		}

		reqToken := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(reqToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized(err.Error()))
			return
		}

		storedToken, err := m.AuthCache.GetTokenCustomer(context.Background(), claims.Email)
		if err != nil || storedToken == nil || *storedToken != reqToken {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized(""))
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("id", claims.ID)
		ctx.Next()
	}
}
