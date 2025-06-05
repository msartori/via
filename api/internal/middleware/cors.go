package middleware

import (
	"net/http"
	"via/internal/logger"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		r = log.WithLogFieldsInRequest(r, "method", r.Method, "uri", r.RequestURI, "proto", r.Proto)
		log.Info(r.Context(), "msg", "CORS middleware invoked")
		//TODO This needs to be configurable
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, bypass-tunnel-reminder")

		// Resolve preflight requests
		// If the request method is OPTIONS, respond with 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
