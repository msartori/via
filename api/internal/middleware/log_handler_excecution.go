package middleware

import (
	"net/http"
	"via/internal/log"
)

func LogHandlerExecution(handlerName string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		logger.Info(r.Context(), "msg", handlerName+".start")
		defer logger.Info(r.Context(), "msg", handlerName+".end")
		next(w, r)
	}
}
