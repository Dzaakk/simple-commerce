package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served.",
		},
		[]string{"method", "route", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDurationSeconds,
		httpRequestsInFlight,
	)
}

// HTTPMiddleware records low-cardinality HTTP metrics for Gin routes.
func HTTPMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isTelemetryEndpoint(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		start := time.Now()
		method := ctx.Request.Method
		inFlightRoute := routePattern(ctx)

		httpRequestsInFlight.WithLabelValues(method, inFlightRoute, "in_flight").Inc()
		defer httpRequestsInFlight.WithLabelValues(method, inFlightRoute, "in_flight").Dec()

		ctx.Next()

		route := routePattern(ctx)
		status := strconv.Itoa(ctx.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(method, route, status).Inc()
		httpRequestDurationSeconds.WithLabelValues(method, route, status).Observe(duration)
	}
}

func routePattern(ctx *gin.Context) string {
	if route := ctx.FullPath(); route != "" {
		return route
	}

	return "unmatched"
}

func isTelemetryEndpoint(path string) bool {
	switch path {
	case "/metrics", "/healthz", "/readyz":
		return true
	default:
		return false
	}
}
