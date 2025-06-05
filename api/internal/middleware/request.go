package middleware

import (
	"context"
	"net/http"
	"time"
	"via/internal/logger"

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
		log := logger.Get()
		r = log.WithLogFieldsInRequest(r, "requestId", reqID)
		log.Info(r.Context(), "msg", "request start")
		// Add Request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)
		// Log when the request finishes
		defer func() {
			duration := time.Since(start)
			log.Info(r.Context(), "msg", "request end", "took", duration.Milliseconds())
		}()
		next.ServeHTTP(w, r)
	})
}

/*
manageRequestId is a helper function to manage the Request ID.
It checks if the Request ID is present in the request header,
and if not, generates a new one.s
It sets the Request ID in the response header.
Uses uuid package to generate a unique ID.
*/
func manageRequestId(w http.ResponseWriter, r *http.Request) string {
	// Get or generate Request ID
	reqID := r.Header.Get(RequestIDHeader)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	w.Header().Set(RequestIDHeader, reqID)
	return reqID
}
