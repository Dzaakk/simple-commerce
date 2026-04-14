package logging

import (
	appLogging "Dzaakk/simple-commerce/package/logging"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger sends basic request/response metadata to Loki.
// It avoids logging request/response bodies to keep payloads safe.
func RequestLogger(client *appLogging.LokiClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		level := "info"
		if status >= http.StatusInternalServerError {
			level = "error"
		} else if status >= http.StatusBadRequest {
			level = "warn"
		}

		fields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  ctx.ClientIP(),
			"user_agent": ctx.Request.UserAgent(),
		}

		if reqID, ok := ctx.Get("request_id"); ok {
			if reqIDStr, ok := reqID.(string); ok && reqIDStr != "" {
				fields["request_id"] = reqIDStr
			}
		} else if reqID := ctx.GetHeader("X-Request-Id"); reqID != "" {
			fields["request_id"] = reqID
		}

		if lastErr := ctx.Errors.Last(); lastErr != nil {
			fields["error"] = lastErr.Err.Error()
		}

		if client == nil {
			return
		}

		pushCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := client.Push(pushCtx, level, "http_request", fields); err != nil {
			log.Printf("loki push failed: %v", err)
		}
	}
}
