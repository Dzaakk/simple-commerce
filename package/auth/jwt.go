package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	model "Dzaakk/simple-commerce/internal/customer/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenPrefix := os.Getenv("TOKEN_PREFIX")
		tokenKey := tokenPrefix + tokenString

		ctxRedis := context.Background()
		customerId, err := redisClient.Get(ctxRedis, tokenKey).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			ctx.Abort()
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					redisClient.Del(ctxRedis, tokenKey)
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
					ctx.Abort()
					return
				}
			}

			// Get customer ID from claims
			if claimCustomerId, ok := claims["customerId"].(string); ok {
				// Verify that Redis customerId matches JWT customerId
				if claimCustomerId != customerId {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token mismatch"})
					ctx.Abort()
					return
				}
				ctx.Set("customerId", claimCustomerId)
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

const TokenExpiration = time.Minute * 30

type TokenGenerator struct {
	redisClient *redis.Client
	jwtSecret   []byte
}

func NewTokenGenerator(redisClient *redis.Client, jwtSecret []byte) *TokenGenerator {
	return &TokenGenerator{
		redisClient: redisClient,
		jwtSecret:   jwtSecret,
	}
}

func (g *TokenGenerator) GenerateToken(customer model.TCustomers) (string, error) {

	now := time.Now()
	expiresAt := now.Add(TokenExpiration)
	claims := jwt.MapClaims{
		"customerId": fmt.Sprintf("%d", customer.Id),
		"exp":        expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(g.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	ctx := context.Background()
	tokenPrefix := os.Getenv("TOKEN_PREFIX")
	tokenKey := tokenPrefix + tokenString

	tokenData := map[string]interface{}{
		"userId":    customer.Id,
		"createdAt": now.Unix(),
		"expiresAt": expiresAt.Unix(),
	}

	tokenDataJson, err := json.Marshal(tokenData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token data: %w", err)
	}

	err = g.redisClient.Set(ctx, tokenKey, tokenDataJson, TokenExpiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return tokenString, nil
}

func (g *TokenGenerator) InvalidateToken(token string) error {
	ctx := context.Background()
	tokenPrefix := os.Getenv("TOKEN_PREFIX")
	tokenKey := tokenPrefix + token

	return g.redisClient.Del(ctx, tokenKey).Err()
}
