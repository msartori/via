package middleware

import (
	"context"
	"net/http"
	"time"
	"via/internal/log"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"
const RequestIDHeader = "x-request-id"

// Middleware to wrap all requests
func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := manageRequestId(w, r)
		logger := log.Get()
		r = logger.WithLogFieldsInRequest(r, "requestId", reqID)
		logger.Info(r.Context(), "msg", "request start")
		// Add Request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)
		// Log when the request finishes
		defer func() {
			duration := time.Since(start)
			logger.Info(r.Context(), "msg", "request end", "took", duration.Milliseconds())
		}()
		next.ServeHTTP(w, r)
	})
}

func manageRequestId(w http.ResponseWriter, r *http.Request) string {
	// Get or generate Request ID
	reqID := r.Header.Get(RequestIDHeader)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	w.Header().Set(RequestIDHeader, reqID)
	return reqID
}
