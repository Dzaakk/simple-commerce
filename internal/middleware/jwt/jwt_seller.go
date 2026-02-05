package middleware

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	authUsecase "Dzaakk/simple-commerce/internal/auth/usecase"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTSellerMiddleware struct {
	TokenUsecase authUsecase.SellerTokenUsecase
}

func NewJWTSellerMiddleware(tokenUsecase authUsecase.SellerTokenUsecase) *JWTSellerMiddleware {
	return &JWTSellerMiddleware{TokenUsecase: tokenUsecase}
}

func (m *JWTSellerMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized("missing or invalid token format"))
			return
		}

		reqToken := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken[*model.SellerToken](reqToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized(err.Error()))
			return
		}

		storedToken, err := m.TokenUsecase.GetToken(context.Background(), claims.Email)
		if err != nil || storedToken == nil || *storedToken != reqToken {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Unauthorized("token mismatch or expired"))
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("id", claims.ID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
