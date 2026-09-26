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

func NewHandler(db *sql.DB, redisClient *redis.Client) *Handler {
	return &Handler{db: db, redis: redisClient}
}

func (h *Handler) Route(router *gin.Engine) {
	router.GET("/healthz", h.Live)
	router.GET("/readyz", h.Ready)
}

func (h *Handler) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": "simple-commerce"})
}

func (h *Handler) Ready(ctx *gin.Context) {
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second)
	defer cancel()
	if err := h.db.PingContext(checkCtx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "postgres": "down"})
		return
	}

	result := gin.H{"status": "ok", "postgres": "ok"}
	if h.redis == nil || h.redis.Ping(checkCtx).Err() != nil {
		result["redis"] = "down; v2 falls back to postgres"
	} else {
		result["redis"] = "ok"
	}
	ctx.JSON(http.StatusOK, result)
}
