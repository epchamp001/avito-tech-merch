package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// HTTPRequestsTotal - счетчик общего количества HTTP-запросов
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration - гистограмма времени выполнения HTTP-запросов
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// HTTPErrorsTotal - счетчик количества ошибок HTTP-запросов
	HTTPErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors.",
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration, HTTPErrorsTotal)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func GinPrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		HTTPRequestsTotal.WithLabelValues(c.Request.Method, endpoint, http.StatusText(status)).Inc()
		HTTPRequestDuration.WithLabelValues(c.Request.Method, endpoint).Observe(duration)
		if status >= 400 {
			HTTPErrorsTotal.WithLabelValues(c.Request.Method, endpoint, http.StatusText(status)).Inc()
		}
	}
}
