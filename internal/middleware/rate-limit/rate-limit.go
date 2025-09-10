package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	redis   *redis.Client
	limiter *rate.Limiter
}

func NewRateLimiter(redis *redis.Client) *RateLimiter {
	return &RateLimiter{
		redis:   redis,
		limiter: rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		clientIP := ctx.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		count, err := rl.redis.Incr(ctx, key).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limit check failed",
			})
			return
		}

		if count == 1 {
			rl.redis.Expire(ctx, key, time.Minute)
		}

		if count > 60 {
			ttl, _ := rl.redis.TTL(ctx, key).Result()
			ctx.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			ctx.Header("X-RateLimit-Limit", "60")
			ctx.Header("X-RateLimit-Remaining", "0")
			ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		ctx.Header("X-RateLimit-Limit", "60")
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(60-count, 10))

		ctx.Next()
	}

}

func (rl *RateLimiter) StrictLimit(requests int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString("id")
		if userID == "" {
			userID = ctx.ClientIP()
		}

		endpoint := ctx.Request.URL.Path
		key := fmt.Sprintf("strict_limit:%s:%s", userID, endpoint)

		now := time.Now().UnixNano()
		windowStart := now - window.Nanoseconds()

		pipe := rl.redis.Pipeline()

		pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

		pipe.ZCard(ctx, key)

		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d", now),
		})

		pipe.Expire(ctx, key, window+time.Minute)

		results, err := pipe.Exec(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limit check failed",
			})
			return
		}

		countCmd := results[1].(*redis.IntCmd)
		count, err := countCmd.Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limit check failed",
			})
			return
		}

		if count >= int64(requests) {
			oldestScore, err := rl.redis.ZRange(ctx, key, 0, 0).Result()
			var resetTime time.Time
			if err != nil && len(oldestScore) > 0 {
				if oldest, err := strconv.ParseInt(oldestScore[0], 10, 64); err == nil {
					resetTime = time.Unix(0, oldest).Add(window)
				}
			}

			if resetTime.IsZero() {
				resetTime = time.Now().Add(window)
			}

			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = int(window.Seconds())
			}

			ctx.Header("Retry-After", strconv.Itoa(retryAfter))
			ctx.Header("X-RateLimit-Limit", strconv.Itoa(requests))
			ctx.Header("X-RateLimit-Remaining", "0")
			ctx.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			ctx.Header("X-RateLimit-Window", window.String())

			rl.redis.ZRem(ctx, key, fmt.Sprintf("%d", now))

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":    "Rate limit exceeded",
				"message":  fmt.Sprintf("Too Many Requests for this endpoint. Try again in %d seconds", retryAfter),
				"endpoint": endpoint,
			})
			return
		}

		remaining := requests - int(count) - 1
		if remaining < 0 {
			remaining = 0
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(requests))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		ctx.Header("X-RateLimit-Window", window.String())

		ctx.Next()
	}
}
