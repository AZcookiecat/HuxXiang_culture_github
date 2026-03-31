package app

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "post_service_requests_total",
			Help: "Total HTTP requests by route and status.",
		}, []string{"route", "method", "status"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "post_service_request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: prometheus.DefBuckets,
		}, []string{"route", "method"}),
	}

	m.requests = registerCounter(m.requests)
	m.latency = registerHistogram(m.latency)
	return m
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := c.Writer.Status()
		m.requests.WithLabelValues(route, c.Request.Method, httpStatus(status)).Inc()
		m.latency.WithLabelValues(route, c.Request.Method).Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) Handler() gin.HandlerFunc {
	inner := promhttp.Handler()
	return func(c *gin.Context) {
		inner.ServeHTTP(c.Writer, c.Request)
	}
}

func registerCounter(counter *prometheus.CounterVec) *prometheus.CounterVec {
	if err := prometheus.Register(counter); err != nil {
		if existing, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return existing.ExistingCollector.(*prometheus.CounterVec)
		}
		panic(err)
	}
	return counter
}

func registerHistogram(histogram *prometheus.HistogramVec) *prometheus.HistogramVec {
	if err := prometheus.Register(histogram); err != nil {
		if existing, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return existing.ExistingCollector.(*prometheus.HistogramVec)
		}
		panic(err)
	}
	return histogram
}

func httpStatus(code int) string {
	return strconv.Itoa(code)
}
