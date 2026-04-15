package logging

import (
	"context"
	"sync"
	"time"
)

type Logger struct {
	client  *LokiClient
	module  string
	service string
}

var (
	defaultClient     *LokiClient
	defaultClientOnce sync.Once
)

func NewLogger(module, service string) *Logger {
	defaultClientOnce.Do(func() {
		defaultClient = NewLokiClientFromEnv()
	})

	return &Logger{
		client:  defaultClient,
		module:  module,
		service: service,
	}
}

func (l *Logger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, "info", message, fields)
}

func (l *Logger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, "warn", message, fields)
}

func (l *Logger) Error(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, "error", message, fields)
}

func (l *Logger) log(ctx context.Context, level, message string, fields map[string]interface{}) {
	if l == nil || l.client == nil {
		return
	}

	payload := map[string]interface{}{
		"module":  l.module,
		"service": l.service,
	}

	for k, v := range fields {
		payload[k] = v
	}

	if ctx != nil {
		if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
			payload["request_id"] = requestID
		}
	}

	pushCtx := context.Background()
	if ctx != nil {
		pushCtx = ctx
	}

	timeoutCtx, cancel := context.WithTimeout(pushCtx, time.Second)
	defer cancel()

	_ = l.client.Push(timeoutCtx, level, message, payload)
}
