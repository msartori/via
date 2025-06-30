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
	"via/internal/response"

	"github.com/go-chi/chi/v5"
)

func New(cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Recover)
	r.Use(middleware.Timeout(time.Duration(cfg.Application.RequestTimeout) * time.Second))
	r.Use(middleware.Request)

	r.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON[any](w, r, response.Response[any]{Data: "ok", Message: "ping status"}, http.StatusOK)
	}))

	// Routes
	r.Get("/guide-to-withdraw/{viaGuideId}", middleware.LogHandlerExecution("handler.GetGuideToWithdraw",
		handler.GetGuideToWithdraw(cfg.Bussiness).ServeHTTP))
	/*
		r.Get("/guide/via-guide-id/{viaGuideId}", middleware.LogHandlerExecution("handler.GetGuideByViaGuideId",
			handler.GetGuideByViaGuideId().ServeHTTP))
	*/
	r.Post("/guide-to-withraw", middleware.LogHandlerExecution("handler.CreateGuideToWidthdraw",
		handler.CreateGuideToWidthdraw(cfg.Bussiness).ServeHTTP))

	r.Get("/monitor/events", middleware.LogHandlerExecution("handler.GetMonitorEvent",
		handler.GetMonitorEvents().ServeHTTP))

	r.Post("/guide/{guideId}/assign", middleware.LogHandlerExecution("handler.AssignGuideToOperator",
		handler.AssignGuideToOperator().ServeHTTP))

	r.Get("/guide/{guideId}/status-options", middleware.LogHandlerExecution("handler.GetGuideStatusOptions",
		handler.GetGuideStatusOptions().ServeHTTP))

	r.Put("/guide/{guideId}/status", middleware.LogHandlerExecution("handler.UpdateGuideStatus",
		handler.UpdateGuideStatus().ServeHTTP))

	r.Get("/auth/login", middleware.LogHandlerExecution("handler.Login",
		handler.Login))

	r.Get("/auth/callback", middleware.LogHandlerExecution("handler.LoginCallback",
		handler.LoginCallback))

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)
		r.Get("/operator/guides", middleware.LogHandlerExecution("handler.GetOperatorGuide",
			handler.GetOperatorGuide().ServeHTTP))
	})

	// Set up dependencies
	via_guide_provider.Set(via_guide_web_provider.New(cfg.GuideWebClient, via_guide_web_provider.HistoricalQueryResponseParser{}))
	guide_provider.Set(guide_ent_provider.New())

	return r
}
