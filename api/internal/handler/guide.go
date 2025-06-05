package handler

import (
	"net/http"
	"time"
	"via/internal/config"
	model_destination "via/internal/model/destination"
	model_guide "via/internal/model/guide"
	"via/internal/secret"
	util_response "via/internal/util/response"
)

func GetGuide() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util_response.ResponderJSON(w, r, model_guide.Guide{ID: 1,
			Test:        secret.ReadSecret(config.Get().Database.PasswordFile) + time.Now().String(),
			Destination: model_destination.Destination{ID: 11}}, http.StatusOK)
	}
}
