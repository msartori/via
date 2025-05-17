package util_response

import (
	"encoding/json"
	"net/http"
)

func ResponderJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ResponderError(w http.ResponseWriter, mensaje string, status int) {
	ResponderJSON(w, map[string]string{"error": mensaje}, status)
}
