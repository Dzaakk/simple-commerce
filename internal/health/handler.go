package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	db    *sql.DB
	redis *redis.Client
}

func NewHandler(db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

func (h *Handler) Route(router *gin.Engine) {
	router.GET("/healthz", h.Live)
	router.GET("/readyz", h.Ready)
}

func (h *Handler) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "simple-commerce",
	})
}

func (h *Handler) Ready(ctx *gin.Context) {
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 1*time.Second)
	defer cancel()

	dependencies := gin.H{}
	ready := true

	if err := h.db.PingContext(checkCtx); err != nil {
		dependencies["postgres"] = "down"
		ready = false
	} else {
		dependencies["postgres"] = "ok"
	}

	if err := h.redis.Ping(checkCtx).Err(); err != nil {
		dependencies["redis"] = "down"
		ready = false
	} else {
		dependencies["redis"] = "ok"
	}

	if !ready {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":       "degraded",
			"service":      "simple-commerce",
			"dependencies": dependencies,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"service":      "simple-commerce",
		"dependencies": dependencies,
	})
}
