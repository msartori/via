package middleware

import (
	"fmt"
	"net/http"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/ratelimit"
	"via/internal/response"
)

type KeyGetter func(r *http.Request) string

type RateLimitMiddleware struct {
	RateLimiter *ratelimit.RateLimiter
	KeyGetter   KeyGetter
}

func NewRateLimitMiddleware(rateLimiters map[string]RateLimitMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rm, ok := rateLimiters[r.URL.Path]
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			key := fmt.Sprintf("%s:%s", rm.RateLimiter.ID, rm.KeyGetter(r))
			allowed, err := rm.RateLimiter.Allow(r.Context(), key)
			if err != nil {
				log.Get().Error(r.Context(), err, "msg", "error checking rate limit",
					"key", key, "id", rm.RateLimiter.ID)
				response.WriteJSON(w, r,
					response.Response[any]{Message: i18n.Get(r, i18n.MsgInternalServerError)}, http.StatusInternalServerError)
				return
			}

			if !allowed {
				log.Get().Warn(r.Context(), "msg", "rate limit",
					"key", key, "id", rm.RateLimiter.ID)
				response.WriteJSON(w, r,
					response.Response[any]{Message: i18n.Get(r, i18n.MsgTooManyRequestsError)}, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
