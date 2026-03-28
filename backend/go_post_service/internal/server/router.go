package server

import (
	"net/http"
	"runtime"

	"huxiang/backend/go_post_service/internal/app"
	"huxiang/backend/go_post_service/internal/community"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg app.Config, metrics *app.Metrics, communityHandler *community.Handler, health func() error) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(app.CORSMiddleware())
	router.Use(app.HTTPSRedirectMiddleware(false))
	router.Use(app.TimeoutMiddleware(cfg.ReadTimeout))
	router.Use(app.NewRateLimiterStore(cfg.RateLimitRPS, cfg.RateLimitBurst).Middleware())
	router.Use(metrics.Middleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "湖湘文化数字化平台 API 接口",
			"version": "go-community-1.0.0",
			"endpoints": gin.H{
				"community": "/api/community",
				"health":    "/health",
				"metrics":   "/metrics",
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		if err := health(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":    "unhealthy",
				"database":  "error",
				"goroutines": runtime.NumGoroutine(),
				"error":     err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"database":  "connected",
			"goroutines": runtime.NumGoroutine(),
		})
	})

	router.GET("/metrics", metrics.Handler())

	api := router.Group("/api")
	communityHandler.Register(api.Group("/community"))

	return router
}
