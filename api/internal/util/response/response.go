package util_response

import (
	"encoding/json"
	"net/http"
	"via/internal/middleware"
)

type Response struct {
	Data      any    `json:"data"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

func ResponderJSON(w http.ResponseWriter, r *http.Request, response Response, status int) {
	if r.Context().Err() != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if requestID, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
		response.RequestID = requestID
	}

	json.NewEncoder(w).Encode(response)
}
