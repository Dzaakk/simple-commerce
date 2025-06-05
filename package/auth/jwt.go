package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"Dzaakk/simple-commerce/internal/customer/model"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const TokenExpiration = time.Minute * 30

func JWTMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractToken(ctx)
		if err != nil {
			response.Unauthorized(err.Error())
			return
		}

		customerId, err := validateSession(ctx, redisClient, tokenString)
		if err != nil {
			handleSessionError(err)
			return
		}

		claims, err := validateToken(tokenString)
		if err != nil {
			response.Unauthorized("Invalid token")
			return
		}
		if isTokenExpired(claims) {
			removeSession(ctx, redisClient, tokenString)
			response.Unauthorized("Token expired")
			return
		}

		if err := validateCustomerId(claims, customerId); err != nil {
			response.Unauthorized(err.Error())
			return
		}

		ctx.Set("customerId", customerId)
		ctx.Next()
	}
}

func extractToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func validateSession(ctx context.Context, redisClient *redis.Client, tokenString string) (string, error) {
	tokenPrefix := os.Getenv("TOKEN_PREFIX")
	tokenKey := tokenPrefix + tokenString

	customerId, err := redisClient.Get(ctx, tokenKey).Result()
	if err != nil {
		return "", err
	}

	return customerId, nil
}

func handleSessionError(err error) {
	if err == redis.Nil {
		response.Unauthorized("Session Expired")
		return
	}

	response.InternalServerError("")
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func isTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		return time.Now().Unix() > int64(exp)
	}

	return true
}

func removeSession(ctx context.Context, redisClient *redis.Client, tokenString string) {
	tokenPrefix := os.Getenv("TOKEN_PREFIX")
	tokenKey := tokenPrefix + tokenString
	redisClient.Del(ctx, tokenKey)
}

func validateCustomerId(claims jwt.MapClaims, redisCustomerId string) error {
	if claimCustomerId, ok := claims["customerId"].(string); ok {
		if claimCustomerId != redisCustomerId {
			return errors.New("token mismatch")
		}
		return nil
	}

	return errors.New("invalid token claims")
}

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
		"customerId": fmt.Sprintf("%d", customer.ID),
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
		"userId":    customer.ID,
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
