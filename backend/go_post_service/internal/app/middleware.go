package app

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rps      rate.Limit
	burst    int
}

func NewRateLimiterStore(rps, burst int) *rateLimiterStore {
	return &rateLimiterStore{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
}

func (s *rateLimiterStore) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := s.getLimiter(ip)
		if !limiter.Allow() {
			AbortError(c, NewError(http.StatusTooManyRequests, "请求过于频繁，请稍后重试", nil))
			return
		}
		c.Next()
	}
}

func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limiter, ok := s.limiters[ip]; ok {
		return limiter
	}
	limiter := rate.NewLimiter(s.rps, s.burst)
	s.limiters[ip] = limiter
	return limiter
}

func RequestLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		logger.Info("http request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("route", route),
			slog.String("client_ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	}
}

func ErrorMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		status, message := ResolveError(err)
		if status >= http.StatusInternalServerError {
			logger.Error("request failed", slog.Any("error", err), slog.Int("status", status))
		}
		Error(c, status, message)
	}
}

func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", slog.Any("error", recovered), slog.String("path", c.Request.URL.Path))
		Error(c, http.StatusInternalServerError, "服务器开小差了，请稍后重试")
		c.Abort()
	})
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Backend-Timeout", timeout.String())
		c.Next()
	}
}

func CORSMiddleware(allowOrigins []string) gin.HandlerFunc {
	allowAll := len(allowOrigins) == 0 || (len(allowOrigins) == 1 && allowOrigins[0] == "*")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if matchedOrigin(origin, allowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func matchedOrigin(origin string, allowOrigins []string) bool {
	for _, allowOrigin := range allowOrigins {
		if strings.EqualFold(strings.TrimSpace(allowOrigin), origin) {
			return true
		}
	}
	return false
}

func HTTPSRedirectMiddleware(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enabled && c.GetHeader("X-Forwarded-Proto") == "http" {
			target := "https://" + c.Request.Host + c.Request.URL.RequestURI()
			c.Redirect(http.StatusPermanentRedirect, target)
			c.Abort()
			return
		}
		c.Next()
	}
}
