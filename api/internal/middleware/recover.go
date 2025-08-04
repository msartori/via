package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"via/internal/log"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Get().Error(r.Context(), fmt.Errorf("%v", rec), "msg", "recovering from panic", "stack", string(debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
