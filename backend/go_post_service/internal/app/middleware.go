package app

import (
	"net/http"
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
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "请求过于频繁，请稍后重试",
			})
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

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(c.Request.Context())
		c.Writer.Header().Set("X-Backend-Timeout", timeout.String())
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
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
