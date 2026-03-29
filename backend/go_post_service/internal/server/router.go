package server

import (
	"log/slog"
	"net/http"
	"runtime"

	"huxiang/backend/go_post_service/internal/app"
	"huxiang/backend/go_post_service/internal/community"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg app.Config, logger *slog.Logger, metrics *app.Metrics, communityHandler *community.Handler, health func() error) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(app.RequestLoggerMiddleware(logger))
	router.Use(app.RecoveryMiddleware(logger))
	router.Use(app.ErrorMiddleware(logger))
	router.Use(app.CORSMiddleware(cfg.CORSAllowOrigins))
	router.Use(app.HTTPSRedirectMiddleware(cfg.EnableHTTPSRedirect))
	router.Use(app.TimeoutMiddleware(cfg.ReadTimeout))
	router.Use(app.NewRateLimiterStore(cfg.RateLimitRPS, cfg.RateLimitBurst).Middleware())
	router.Use(metrics.Middleware())

	router.GET("/", func(c *gin.Context) {
		app.Success(c, http.StatusOK, gin.H{
			"message": "湖湘文化数字化平台 API 接口",
			"version": "go-community-1.1.0",
			"endpoints": gin.H{
				"community": "/api/community",
				"health":    "/health",
				"metrics":   "/metrics",
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		if err := health(); err != nil {
			app.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		app.Success(c, http.StatusOK, gin.H{
			"status":     "healthy",
			"database":   "connected",
			"goroutines": runtime.NumGoroutine(),
		})
	})

	router.GET("/metrics", metrics.Handler())

	api := router.Group("/api")
	communityHandler.Register(api.Group("/community"))
	return router
}
