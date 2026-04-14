package requestid

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const headerRequestID = "X-Request-Id"

// RequestID ensures every request has a request id.
// It reuses incoming X-Request-Id if provided, otherwise generates one.
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.GetHeader(headerRequestID)
		if reqID == "" {
			reqID = generateRequestID()
		}

		ctx.Set("request_id", reqID)
		ctx.Request.Header.Set(headerRequestID, reqID)
		ctx.Header(headerRequestID, reqID)

		ctx.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b)
}
