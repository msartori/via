package router

import (
	"net/http"
	"time"

	"via/internal/config"
	"via/internal/handler"
	"via/internal/middleware"

	guide_provider "via/internal/provider/guide"
	guide_web_provider "via/internal/provider/guide/web"

	"github.com/go-chi/chi/v5"
)

func New(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.Timeout(time.Duration(cfg.Application.RequestTimeout) * time.Second))
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Request)

	// Routes
	r.Get("/guide/{id}", handler.GetGuide(cfg.Bussiness).ServeHTTP)

	// Set up dependencies
	guideProvider := guide_web_provider.New(cfg.GuideWebClient, guide_web_provider.HistoricalQueryResponseParser{})
	guide_provider.Set(guideProvider)

	return r
}
