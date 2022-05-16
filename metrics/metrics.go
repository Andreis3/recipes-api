package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of incoming requests",
	},
	[]string{"path", "method"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time spent serving HTTP requests",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	},
	[]string{"path", "method"},
)

var statusCodes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_status_codes_total",
		Help: "Number of HTTP status codes",
	},
	[]string{"path", "method", "code"},
)

func PrometheusMiddleware(urlsExcludedFromMonitoring []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		incomingRequestUrl := c.Request.URL.Path
		for _, v := range urlsExcludedFromMonitoring {
			if incomingRequestUrl == v {
				return
			}
		}

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(c.Request.URL.Path, c.Request.Method))
		totalRequests.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()
		statusCodes.WithLabelValues(c.Request.URL.Path, c.Request.Method, fmt.Sprintf("%d", c.Writer.Status())).Inc()

		timer.ObserveDuration()
	}
}

func RegisterPrometheus() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(statusCodes)
}
