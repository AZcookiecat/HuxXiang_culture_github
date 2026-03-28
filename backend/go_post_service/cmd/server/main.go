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

	dbs, err := app.OpenDatabases(cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer dbs.Writer.Close()
	defer dbs.Reader.Close()

	cache := app.NewInMemoryCache()
	events := app.NewEventBus(cache)
	service := community.NewService(community.NewMySQLRepository(dbs.Writer, dbs.Reader), cache, events, cfg.CacheTTL)
	handler := community.NewHandler(service, cfg.JWTSecret)
	metrics := app.NewMetrics()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go events.Run(ctx)

	router := server.NewRouter(cfg, metrics, handler, func() error {
		return service.Health(context.Background())
	})

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Printf("go post service listening on %s", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpServer.Shutdown(shutdownCtx)
}
