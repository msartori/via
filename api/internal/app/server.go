package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"via/internal/config"
	"via/internal/log"
	"via/internal/router"
)

func StartServer(cfg config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := router.New(cfg)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Application.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	logger := log.Get()
	go func() {
		logger.Info(ctx, "msg", "Server started", "port", cfg.Application.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, err, "msg", "Server failed")
		}
	}()

	<-ctx.Done()
	logger.Info(context.Background(), "msg", "Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, err, "msg", "Graceful shutdown failed")
	} else {
		logger.Info(ctx, "msg", "Server stopped gracefully")
	}
}
