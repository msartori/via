package server

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"via/internal/log"
)

type ServerConfig struct {
	Enabled      bool `env:"ENABLED" envDefault:"true"     json:"enabled"`
	Port         int  `env:"PORT" envDefault:"8080"        json:"port"`
	ReadTimeout  int  `env:"READ_TIMEOUT" envDefault:"15"  json:"readTimeout"`
	WriteTimeout int  `env:"WRITE_TIMEOUT" envDefault:"15" json:"writeTimeout"`
	IdleTimeout  int  `env:"IDLE_TIMEOUT" envDefault:"30"  json:"idleTimeout"`
}

func Create(cfg ServerConfig, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}
}

func Start(servers []*http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.Get()

	// Launch each in a go routine
	for _, srv := range servers {
		go func(s *http.Server) {
			logger.Info(ctx, "msg", "Server started", "port", s.Addr)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal(ctx, err, "msg", "Server failed", "port", s.Addr)
			}
		}(srv)
	}

	// Wait for shutdown
	<-ctx.Done()
	logger.Info(context.Background(), "msg", "Shutdown signal received")

	// Shutdown each in parallel
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, srv := range servers {
		wg.Add(1)
		go func(s *http.Server) {
			defer wg.Done()
			if err := s.Shutdown(shutdownCtx); err != nil {
				logger.Error(ctx, err, "msg", "Graceful shutdown failed", "port", s.Addr)
			} else {
				logger.Info(ctx, "msg", "Server stopped gracefully", "port", s.Addr)
			}
		}(srv)
	}
	wg.Wait()
}

/*

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
*/
