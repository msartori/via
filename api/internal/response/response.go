package response

import (
	"encoding/json"
	"net/http"
	"strings"
	biz_language "via/internal/biz/language"
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

func GetLanguage(r *http.Request) string {
	langHeader := r.Header.Get("Accept-Language")
	if langHeader == "" {
		return biz_language.DEFAULT //default language
	}

	languages := strings.Split(langHeader, ",")

	if len(languages) == 0 {
		return biz_language.DEFAULT
	}

	// Get the first and clean it
	primary := strings.SplitN(strings.TrimSpace(languages[0]), ";", 2)[0]

	if primary == "" {
		return biz_language.DEFAULT
	}

	return primary
}
