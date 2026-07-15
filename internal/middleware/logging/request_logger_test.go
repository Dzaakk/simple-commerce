package logging

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type capturedLog struct {
	level   string
	message string
	fields  map[string]interface{}
}

type captureLogger struct {
	entries []capturedLog
}

func (l *captureLogger) Info(_ context.Context, message string, fields map[string]interface{}) {
	l.capture("info", message, fields)
}

func (l *captureLogger) Warn(_ context.Context, message string, fields map[string]interface{}) {
	l.capture("warn", message, fields)
}

func (l *captureLogger) Error(_ context.Context, message string, fields map[string]interface{}) {
	l.capture("error", message, fields)
}

func (l *captureLogger) capture(level, message string, fields map[string]interface{}) {
	l.entries = append(l.entries, capturedLog{level: level, message: message, fields: fields})
}

func TestRequestLoggerEmitsExpectedLevelAndFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		status    int
		withError bool
		wantLevel string
	}{
		{name: "success", status: http.StatusOK, wantLevel: "info"},
		{name: "client error", status: http.StatusBadRequest, wantLevel: "warn"},
		{name: "server error", status: http.StatusInternalServerError, withError: true, wantLevel: "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &captureLogger{}
			router := gin.New()
			router.Use(func(ctx *gin.Context) {
				ctx.Set("request_id", "req-123")
				ctx.Next()
			})
			router.Use(RequestLogger(logger))
			router.GET("/products/:id", func(ctx *gin.Context) {
				if tt.withError {
					_ = ctx.Error(errors.New("repository unavailable"))
				}
				ctx.Status(tt.status)
			})

			request := httptest.NewRequest(http.MethodGet, "/products/product-1", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if len(logger.entries) != 1 {
				t.Fatalf("log entries = %d, want 1", len(logger.entries))
			}

			entry := logger.entries[0]
			if entry.level != tt.wantLevel {
				t.Fatalf("level = %q, want %q", entry.level, tt.wantLevel)
			}
			if entry.message != "http_request" {
				t.Fatalf("message = %q, want http_request", entry.message)
			}
			if entry.fields["route"] != "/products/:id" {
				t.Fatalf("route = %#v, want /products/:id", entry.fields["route"])
			}
			if entry.fields["request_id"] != "req-123" {
				t.Fatalf("request_id = %#v, want req-123", entry.fields["request_id"])
			}
			if entry.fields["status"] != tt.status {
				t.Fatalf("status = %#v, want %d", entry.fields["status"], tt.status)
			}
			if tt.withError && entry.fields["error"] != "repository unavailable" {
				t.Fatalf("error = %#v, want repository unavailable", entry.fields["error"])
			}
		})
	}
}

func TestRequestLoggerSkipsTelemetryEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := &captureLogger{}
	requestLogger := RequestLogger(logger)

	for _, path := range []string{"/healthz", "/readyz", "/metrics"} {
		router := gin.New()
		router.Use(requestLogger)
		router.GET(path, func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}

	if len(logger.entries) != 0 {
		t.Fatalf("log entries = %d, want 0", len(logger.entries))
	}
}
