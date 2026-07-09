package middleware

import (
	"net/http"
	"os"
	"strings"

	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := bearerToken(ctx.GetHeader("Authorization"))
		if tokenString == "" {
			ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
			ctx.Abort()
			return
		}

		claims, err := parseAccessToken(tokenString)
		if err != nil || claims.UserID == "" || claims.UserType == "" {
			ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
			ctx.Abort()
			return
		}

		ctx.Set("id", claims.UserID)
		ctx.Set("user_type", claims.UserType)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

func RequireUserType(userType constant.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetString("user_type") != string(userType) {
			ctx.Error(response.NewAppError(http.StatusForbidden, "forbidden"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func parseAccessToken(tokenString string) (*accessTokenClaims, error) {
	claims := &accessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
