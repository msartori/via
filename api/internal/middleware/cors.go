package middleware

import (
	"net/http"
	"via/internal/config"
	"via/internal/log"
)

func CORS(next http.Handler, cfg config.CORS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		r = logger.WithLogFieldsInRequest(r, "method", r.Method, "uri", r.RequestURI, "proto", r.Proto)
		logger.Info(r.Context(), "msg", "CORS middleware invoked")
		//TODO This needs to be configurable
		w.Header().Set("Access-Control-Allow-Origin", cfg.Origins)
		w.Header().Set("Access-Control-Allow-Methods", cfg.Methods)
		w.Header().Set("Access-Control-Allow-Headers", cfg.Headers)

		// Resolve preflight requests
		// If the request method is OPTIONS, respond with 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
