package middleware

import (
	"net/http"
	"strings"
	"via/internal/log"
)

type CORSCfg struct {
	Enabled bool   `env:"ENABLED" envDefault:"true" json:"enabled"`
	Origins string `env:"ORIGINS" envDefault:"*" json:"origins"`
	Methods string `env:"METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS" json:"methods"`
	Headers string `env:"HEADERS" envDefault:"Content-Type,Authorization,bypass-tunnel-reminder,Accept-Language" json:"headers"`
}

func CORS(cfg CORSCfg) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.Get()
			r = logger.WithLogFieldsInRequest(r, "method", r.Method, "uri", r.RequestURI, "proto", r.Proto)
			logger.Info(r.Context(), "msg", "CORS middleware invoked")
			origin := r.Header.Get("Origin")
			if cfg.Origins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", cfg.Origins)
			} else {
				origins := strings.SplitSeq(cfg.Origins, ",")
				for o := range origins {
					//	log.Get().Info(r.Context(), "origin", origins, "o", o)
					if o == origin {
						w.Header().Set("Access-Control-Allow-Origin", o)
						break
					}
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", cfg.Methods)
			w.Header().Set("Access-Control-Allow-Headers", cfg.Headers)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Resolve preflight requests
			// If the request method is OPTIONS, respond with 200 OK
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
