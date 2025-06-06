package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"via/internal/config"
	"via/internal/handler"
	"via/internal/log"
	app_log "via/internal/log/app"
	"via/internal/middleware"
)

func main() {
	// Initialize the logger
	log.Set(app_log.New(config.Get().Log))
	logger := log.Get()

	mux := http.NewServeMux()
	mux.HandleFunc("/guide", handler.GetGuide())
	ctx := context.Background()
	cfg := config.Get()
	logger.Info(ctx, "config", cfg)
	logger.Info(ctx, "msg", "Server Started")
	logger.Fatal(ctx, http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port),
		middleware.CORS(middleware.Recover(middleware.Request(middleware.Timeout(5*time.Second)(mux))), cfg.CORS)))
}
