package logging

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Info(context.Context, string, map[string]interface{})
	Warn(context.Context, string, map[string]interface{})
	Error(context.Context, string, map[string]interface{})
}

// RequestLogger emits request metadata to standard output through the supplied
// structured logger. Log collection and delivery are handled outside the API
// request path by Promtail.
func RequestLogger(logger Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isTelemetryEndpoint(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		start := time.Now()
		ctx.Next()

		fields := map[string]interface{}{
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"status":     ctx.Writer.Status(),
			"latency_ms": time.Since(start).Milliseconds(),
			"client_ip":  ctx.ClientIP(),
			"user_agent": ctx.Request.UserAgent(),
		}

		if route := ctx.FullPath(); route != "" {
			fields["route"] = route
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

		if logger == nil {
			return
		}

		switch status := ctx.Writer.Status(); {
		case status >= http.StatusInternalServerError:
			logger.Error(ctx.Request.Context(), "http_request", fields)
		case status >= http.StatusBadRequest:
			logger.Warn(ctx.Request.Context(), "http_request", fields)
		default:
			logger.Info(ctx.Request.Context(), "http_request", fields)
		}
	}
}

func isTelemetryEndpoint(path string) bool {
	switch path {
	case "/metrics", "/healthz", "/readyz":
		return true
	default:
		return false
	}
}
