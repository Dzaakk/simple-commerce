package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestLoggerEmitsStructuredJSON(t *testing.T) {
	var output bytes.Buffer
	base := slog.New(slog.NewJSONHandler(&output, nil))
	logger := newLogger(base, "user", "customer_service")
	ctx := context.WithValue(context.Background(), "request_id", "req-123")

	logger.Error(ctx, "customer_create_failed", map[string]interface{}{
		"customer_id": "customer-1",
		"attempt":     2,
	})

	var entry map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("decode log entry: %v", err)
	}

	assertLogField(t, entry, "level", "ERROR")
	assertLogField(t, entry, "msg", "customer_create_failed")
	assertLogField(t, entry, "module", "user")
	assertLogField(t, entry, "service", "customer_service")
	assertLogField(t, entry, "request_id", "req-123")
	assertLogField(t, entry, "customer_id", "customer-1")

	if got := entry["attempt"]; got != float64(2) {
		t.Fatalf("attempt = %#v, want 2", got)
	}
}

func TestLoggerAcceptsNilContext(t *testing.T) {
	var output bytes.Buffer
	base := slog.New(slog.NewJSONHandler(&output, nil))
	logger := newLogger(base, "email", "publisher")

	logger.Info(nil, "message_queued", nil)

	if output.Len() == 0 {
		t.Fatal("expected a JSON log entry")
	}
}

func assertLogField(t *testing.T, entry map[string]interface{}, key string, want interface{}) {
	t.Helper()
	if got := entry[key]; got != want {
		t.Fatalf("%s = %#v, want %#v", key, got, want)
	}
}
