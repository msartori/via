package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"via/internal/config"
	"via/internal/handler"
	"via/internal/logger"
	"via/internal/middleware"
)

func main() {
	log := logger.Get()
	mux := http.NewServeMux()
	mux.HandleFunc("/guide", handler.GetGuide())
	ctx := context.Background()
	log.Info(ctx, "config", config.Get())
	log.Info(ctx, "msg", "Server Started")
	log.Fatal(ctx, http.ListenAndServe(fmt.Sprintf(":%d", config.Get().Application.Port),
		middleware.CORS(middleware.Recover(middleware.Request(middleware.Timeout(5*time.Second)(mux))))))
}
