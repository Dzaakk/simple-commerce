package logging

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	logger  *slog.Logger
	module  string
	service string
}

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func NewLogger(module, service string) *Logger {
	return newLogger(defaultLogger, module, service)
}

func newLogger(logger *slog.Logger, module, service string) *Logger {
	return &Logger{
		logger:  logger,
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
	if l == nil || l.logger == nil {
		return
	}

	attrs := make([]slog.Attr, 0, len(fields)+3)
	attrs = append(attrs,
		slog.String("module", l.module),
		slog.String("service", l.service),
	)

	if ctx != nil {
		if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
			attrs = append(attrs, slog.String("request_id", requestID))
		}
	} else {
		ctx = context.Background()
	}

	for key, value := range fields {
		attrs = append(attrs, slog.Any(key, value))
	}

	l.logger.LogAttrs(ctx, parseLevel(level), message, attrs...)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
