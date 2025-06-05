package util_response

import (
	"encoding/json"
	"net/http"
)

func ResponderJSON(w http.ResponseWriter, r *http.Request, data any, status int) {

	if r.Context().Err() != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ResponderError(w http.ResponseWriter, r *http.Request, mensaje string, status int) {
	ResponderJSON(w, r, map[string]string{"error": mensaje}, status)
}
