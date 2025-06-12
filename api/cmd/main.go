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
	guide_provider "via/internal/provider/guide"
	guide_web_provider "via/internal/provider/guide/web"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Get()
	// Initialize the logger
	log.Set(app_log.New(cfg.Log))
	logger := log.Get()

	guideProvider := guide_web_provider.New(cfg.GuideWebClient)
	guide_provider.Set(guideProvider)

	router := chi.NewRouter()
	router.Get("/guide/{id}", handler.GetGuide(cfg.Bussiness).ServeHTTP)
	/*
		router.Get("/guide/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			fmt.Fprintf(w, "Guide ID: %s", id)
		})*/
	ctx := context.Background()

	logger.Info(ctx, "config", cfg)
	logger.Info(ctx, "msg", "Server Started")

	logger.Fatal(ctx, http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port),
		middleware.CORS(middleware.Recover(middleware.Request(middleware.Timeout(time.Duration(cfg.Application.RequestTimeout)*time.Second)(router))), cfg.CORS)))

	//logger.Fatal(ctx, http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), router))
}
