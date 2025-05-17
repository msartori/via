package handler_guide

import (
	"net/http"
	model_destination "via/internal/model/destination"
	model_guide "via/internal/model/guide"
	util_response "via/internal/util/response"
)

func GetGuide() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util_response.ResponderJSON(w, model_guide.Guide{ID: 1, Destination: model_destination.Destination{ID: 11}}, http.StatusOK)
	}
}
