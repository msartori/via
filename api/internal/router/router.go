package router

import (
	"net/http"
	"time"
	"via/internal/config"
	"via/internal/ds"
	"via/internal/global"
	"via/internal/handler"
	"via/internal/middleware"
	"via/internal/ratelimit"
	"via/internal/response"
	"via/internal/sse"

	"github.com/go-chi/chi/v5"
)

var rateFunctionByIP = func(r *http.Request) string {
	return r.RemoteAddr
}

func NewRest(cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Recover)
	r.Use(middleware.Timeout(time.Duration(cfg.Application.RequestTimeout) * time.Second))
	r.Use(middleware.Request)
	r.Use(middleware.NewRateLimitMiddleware(
		map[string]middleware.RateLimitMiddleware{
			"/auth/login": {
				RateLimiter: ratelimit.New("Login", 6, 1*time.Minute, ds.Get()),
				KeyGetter:   rateFunctionByIP,
			},
		},
	))

	r.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, r, response.Response[any]{Data: "ok", Message: "ping status"}, http.StatusOK)
	}))

	// Routes
	r.Get("/guide-to-withdraw/{viaGuideId}", middleware.LogHandlerExecution("handler.GetGuideToWithdraw",
		handler.GetGuideToWithdraw(cfg.Bussiness).ServeHTTP))

	r.Post("/guide-to-withdraw", middleware.LogHandlerExecution("handler.CreateGuideToWidthdraw",
		handler.CreateGuideToWidthdraw(cfg.Bussiness).ServeHTTP))

	r.Get("/auth/login", middleware.LogHandlerExecution("handler.Login",
		handler.Login().ServeHTTP))

	r.Get("/auth/callback", middleware.LogHandlerExecution("handler.LoginCallback",
		handler.LoginCallback(cfg.OAuth).ServeHTTP))

	r.Post("/auth/logout", middleware.LogHandlerExecution("handler.LogOut",
		handler.LogOut(cfg.OAuth).ServeHTTP))

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.OAuth))

		r.Post("/guide/{guideId}/assign", middleware.LogHandlerExecution("handler.AssignGuideToOperator",
			handler.AssignGuideToOperator().ServeHTTP))

		r.Get("/guide/{guideId}/status-options", middleware.LogHandlerExecution("handler.GetGuideStatusOptions",
			handler.GetGuideStatusOptions().ServeHTTP))

		r.Put("/guide/{guideId}/status", middleware.LogHandlerExecution("handler.UpdateGuideStatus",
			handler.UpdateGuideStatus().ServeHTTP))
	})

	return r
}

func NewSSE(cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Recover)
	r.Use(middleware.Request)
	r.Use(middleware.NewRateLimitMiddleware(
		map[string]middleware.RateLimitMiddleware{
			"/auth/login": {
				RateLimiter: ratelimit.New("Login", 6, 1*time.Minute, ds.Get()),
				KeyGetter:   rateFunctionByIP,
			},
			"/operator/guides": {
				RateLimiter: ratelimit.New("OperatorGuides", 6, 1*time.Minute, ds.Get()),
				KeyGetter:   rateFunctionByIP,
			},
		},
	))
	// Routes
	r.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, r, response.Response[any]{Data: "ok", Message: "ping status"}, http.StatusOK)
	}))

	r.Get("/monitor/events", middleware.LogHandlerExecution("handler.GetMonitorEvents",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sse.HandleSSE(w, r, handler.GetMonitorEvents, global.NewGuideChannel, global.GuideStatusChangeChannel)
		})))

	r.Get("/auth/login", middleware.LogHandlerExecution("handler.Login",
		handler.Login().ServeHTTP))

	r.Get("/auth/callback", middleware.LogHandlerExecution("handler.LoginCallback",
		handler.LoginCallback(cfg.OAuth).ServeHTTP))

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.OAuth))

		r.Get("/operator/guides", middleware.LogHandlerExecution("handler.GetOperatorGuides",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sse.HandleSSE(w, r, handler.GetOperatorGuide, global.NewGuideChannel, global.GuideAssignmentChannel, global.GuideStatusChangeChannel)
			})))
	})
	return r
}
