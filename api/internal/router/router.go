package router

import (
	"net/http"
	"time"
	"via/internal/config"
	"via/internal/handler"
	"via/internal/middleware"
	guide_provider "via/internal/provider/guide"
	guide_ent_provider "via/internal/provider/guide/ent"
	via_guide_provider "via/internal/provider/via/guide"
	via_guide_web_provider "via/internal/provider/via/guide/web"

	"github.com/go-chi/chi/v5"
)

func New(cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Recover)
	r.Use(middleware.Timeout(time.Duration(cfg.Application.RequestTimeout) * time.Second))
	r.Use(middleware.Request)

	// Routes
	r.Get("/guide-to-withdraw/{viaGuideId}", middleware.LogHandlerExecution("handler.GetGuideToWithdraw",
		handler.GetGuideToWithdraw(cfg.Bussiness).ServeHTTP))
	r.Get("/guide/via-guide-id/{viaGuideId}", middleware.LogHandlerExecution("handler.GetGuideByViaGuideId",
		handler.GetGuideByViaGuideId().ServeHTTP))
	r.Post("/guide-to-withraw", middleware.LogHandlerExecution("handler.CreateGuideToWidthdraw",
		handler.CreateGuideToWidthdraw(cfg.Bussiness).ServeHTTP))

	// Set up dependencies
	via_guide_provider.Set(via_guide_web_provider.New(cfg.GuideWebClient, via_guide_web_provider.HistoricalQueryResponseParser{}))
	guide_provider.Set(guide_ent_provider.New())

	return r
}
