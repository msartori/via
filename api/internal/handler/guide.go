package handler

import (
	"net/http"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	via_guide_provider "via/internal/provider/via/guide"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

func GetGuideByViaGuideId() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		res := response.Response[model.Guide]{}
		viaGuideId := chi.URLParam(r, "viaGuideId")

		if ok := isValidViaGuideId(w, r, viaGuideId); !ok {
			return
		}

		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "via_guide_id", viaGuideId)

		guide, err := guide_provider.Get().GetGuideByViaGuideId(r.Context(), viaGuideId)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		if ok := isGuideNotFound(w, r, guide.ViaGuideID); !ok {
			return
		}

		res.Data = guide
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type CreateGuideToWidthdrawInput struct {
	ViaGuideId string `json:"viaGuideId"`
}

type CreateGuideToWidthdrawOutput struct {
	ID int `json:"id"`
}

func CreateGuideToWidthdraw() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[*CreateGuideToWidthdrawOutput]{}
		var input CreateGuideToWidthdrawInput

		if ok := getJsonBody(w, r, &input); !ok {
			return
		}

		if ok := isValidViaGuideId(w, r, input.ViaGuideId); !ok {
			return
		}
		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "via_guide_id", input.ViaGuideId)

		viaGuide, err := via_guide_provider.Get().GetGuide(r.Context(), input.ViaGuideId)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		data := &CreateGuideToWidthdrawOutput{}
		res.Data = data

		if ok := isGuideNotFound(w, r, viaGuide.ID); !ok {
			return
		}

		id, err := guide_provider.Get().CreateGuide(r.Context(), viaGuide)

		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to create guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}
		data.ID = id
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}
