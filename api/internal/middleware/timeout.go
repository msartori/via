package middleware

import (
	"context"
	"net/http"
	"time"
	"via/internal/logger"
)

const HttpStatusClientCloseRequest = 499

// timeoutWriter captures response for potential delayed write
type timeoutWriter struct {
	http.ResponseWriter
	header     http.Header
	statusCode int
	body       []byte
	written    bool
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.body = append(tw.body, b...)
	tw.written = true
	return len(b), nil
}

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.statusCode = statusCode
}

func (tw *timeoutWriter) copy(r *http.Request) {
	// Only write if not already written by middleware
	if !tw.written {
		return
	}
	h := tw.ResponseWriter.Header()
	for k, v := range tw.header {
		h[k] = v
	}
	if tw.statusCode != 0 {
		tw.ResponseWriter.WriteHeader(tw.statusCode)
		logger.Get().WithLogFieldsInRequest(r, "status", tw.statusCode)
	}
	tw.ResponseWriter.Write(tw.body)
}

// Middleware to apply a global timeout
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			//release the context when done
			defer cancel()
			// Replace the request context with the new context that has a timeout
			r = r.WithContext(ctx)
			// Create a channel to signal when the handler is done
			done := make(chan struct{})

			// Create a ResponseWriter wrapper to avoid writing twice
			tw := &timeoutWriter{ResponseWriter: w, header: http.Header{}}

			// handle the request in a goroutine
			// This allows us to wait for the handler to finish or timeout
			// without blocking the main goroutine
			go func() {
				next.ServeHTTP(tw, r)
				close(done)
			}()
			statusCode := http.StatusOK
			msg := ""
			var err error
			log := logger.Get()
			select {
			case <-ctx.Done(): //timeout or cancellation
				switch err = ctx.Err(); err {
				case context.DeadlineExceeded:
					msg = "request timeout"
					statusCode = http.StatusGatewayTimeout
				case context.Canceled:
					msg = "request canceled by client"
					statusCode = HttpStatusClientCloseRequest
				default:
					msg = "unexpected context error"
					statusCode = http.StatusServiceUnavailable
				}
				log.WithLogFieldsInRequest(r, "status", statusCode)
				log.Error(r.Context(), err, "msg", msg)
				w.WriteHeader(statusCode)
			case <-done: // handler completed successfully
				tw.copy(r)
			}
		})
	}
}
