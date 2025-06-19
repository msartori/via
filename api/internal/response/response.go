package response

import (
	"encoding/json"
	"net/http"
	"via/internal/global"
)

type Response[T any] struct {
	Data      T      `json:"data"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

func WriteJSON[T any](w http.ResponseWriter, r *http.Request, res Response[T], status int) {
	if r.Context().Err() != nil {
		return
	}
	writeJSON(w, r, res, status)
}

func WriteJSONError[T any](w http.ResponseWriter, r *http.Request, res Response[T], status int) {
	writeJSON(w, r, res, status)
}

func writeJSON[T any](w http.ResponseWriter, r *http.Request, res Response[T], status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if requestID, ok := r.Context().Value(global.RequestIDKey).(string); ok {
		res.RequestID = requestID
	}
	json.NewEncoder(w).Encode(res)
}
