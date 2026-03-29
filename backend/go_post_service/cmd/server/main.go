package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"huxiang/backend/go_post_service/internal/app"
	"huxiang/backend/go_post_service/internal/community"
	"huxiang/backend/go_post_service/internal/server"
)

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := app.NewLogger(cfg)

	dbs, err := app.OpenDatabases(cfg)
	if err != nil {
		logger.Error("open database failed", "error", err)
		os.Exit(1)
	}
	defer dbs.Writer.Close()
	defer dbs.Reader.Close()

	cache := app.NewInMemoryCache()
	events := app.NewEventBus(cache)
	metrics := app.NewMetrics()
	repo := community.NewMySQLRepository(dbs.Writer, dbs.Reader)
	service := community.NewService(repo, cache, events, cfg.CacheTTL)
	handler := community.NewHandler(service, cfg.JWTSecret)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go events.Run(ctx)

	router := server.NewRouter(cfg, logger, metrics, handler, func() error {
		return service.Health(context.Background())
	})

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		logger.Info("server listening", "addr", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server start failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpServer.Shutdown(shutdownCtx)
}
