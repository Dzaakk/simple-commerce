package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type LokiClient struct {
	BaseURL string
	Labels  map[string]string
	Client  *http.Client
}

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type lokiLine struct {
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Time    string                 `json:"time"`
}

func NewLokiClient(baseURL string, labels map[string]string) *LokiClient {
	if baseURL == "" {
		baseURL = "http://localhost:3100"
	}
	if labels == nil {
		labels = map[string]string{}
	}

	if _, ok := labels["job"]; !ok {
		labels["job"] = "simple-commerce"
	}

	return &LokiClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Labels:  labels,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func NewLokiClientFromEnv() *LokiClient {
	baseURL := os.Getenv("LOKI_URL")
	if strings.TrimSpace(baseURL) == "" {
		return nil
	}

	labels := map[string]string{}

	if app := os.Getenv("LOKI_APP"); app != "" {
		labels["app"] = sanitizeLabel(app)
	}
	if env := os.Getenv("LOKI_ENV"); env != "" {
		labels["env"] = sanitizeLabel(env)
	}

	return NewLokiClient(baseURL, labels)
}

func (c *LokiClient) Push(ctx context.Context, level, message string, fields map[string]interface{}) error {
	if c == nil {
		return fmt.Errorf("loki client is nil")
	}

	line := lokiLine{
		Level:   level,
		Message: message,
		Fields:  fields,
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
	}

	lineBytes, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("marshal log line: %w", err)
	}

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: c.Labels,
				Values: [][]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), string(lineBytes)},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal loki payload: %w", err)
	}

	endpoint := c.BaseURL + "/loki/api/v1/push"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create loki request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("send loki request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki push failed: status %d", resp.StatusCode)
	}

	return nil
}

func sanitizeLabel(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, " ", "_")
	val = strings.ToLower(val)
	return val
}
