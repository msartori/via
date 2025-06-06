package middleware

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"sync"
	"time"
	"via/internal/log"
)

const HttpStatusClientCloseRequest = 499

// timeoutWriter captures response for potential delayed write
type timeoutWriter struct {
	http.ResponseWriter
	header     http.Header
	statusCode int
	body       []byte
	written    bool
	timedOut   bool
	mu         sync.Mutex
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return 0, errors.New("write after timeout")
	}
	tw.body = append(tw.body, b...)
	tw.written = true
	return len(b), nil
}

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.statusCode = statusCode
}

func (tw *timeoutWriter) copy(r *http.Request) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut || !tw.written {
		return
	}
	h := tw.ResponseWriter.Header()
	maps.Copy(h, tw.header)
	if tw.statusCode != 0 {
		tw.ResponseWriter.WriteHeader(tw.statusCode)
		log.Get().WithLogFieldsInRequest(r, "status", tw.statusCode)
	}
	tw.ResponseWriter.Write(tw.body)
}

func (tw *timeoutWriter) markTimedOut() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.timedOut = true
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
			logger := log.Get()
			select {
			case <-ctx.Done(): //timeout or cancellation
				tw.markTimedOut()
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
				logger.WithLogFieldsInRequest(r, "status", statusCode)
				logger.Error(r.Context(), err, "msg", msg)
				w.WriteHeader(statusCode)
			case <-done: // handler completed successfully
				tw.copy(r)
			}
		})
	}
}
