package auth

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
		tokenKey := os.Getenv("TOKEN_KEY")
		ctxRedis := context.Background()
		var cursor uint64
		// Retrieve the token associated with the customer ID from Redis
		storedToken, _, err := redisClient.Scan(ctxRedis, cursor, tokenKey, 0).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			ctx.Abort()
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			ctx.Abort()
			return
		}

		redisToken := strings.Split(storedToken[0], ":")
		// Validate the token retrieved from Redis
		token, err := jwt.Parse(redisToken[2], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token Expired"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if customerId, ok := claims["customerId"].(string); ok {
				ctx.Set("customerId", customerId) // Store the user ID in the context
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func TokenJWTGenerator(redisClient *redis.Client, customer model.TCustomers) (string, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	claims := jwt.MapClaims{
		"customerId": fmt.Sprintf("%d", customer.Id),
		"exp":        time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	tokenPattern := os.Getenv("TOKEN_PATTERN")
	ctx := context.Background()
	tokenKey := tokenPattern + tokenString

	err = redisClient.Set(ctx, tokenKey, customer.Id, time.Minute*5).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
