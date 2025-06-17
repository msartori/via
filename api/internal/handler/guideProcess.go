package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"via/internal/i18n"
	"via/internal/log"
	guide_provider "via/internal/provider/guide"
	guide_process_provider "via/internal/provider/guide/process"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

func GetGuideProcessByCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		logger.Info(r.Context(), "msg", "handler.GetGuideProcessByCode_start")
		defer logger.Info(r.Context(), "msg", "handler.GetGuideByCode_end")

		res := response.Response{}
		code := chi.URLParam(r, "code")

		if code == "" {
			logger.Warn(r.Context(), "msg", "missing guide code")
			res.Message = i18n.Get(r, i18n.MsgGuideRequired)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		if match, _ := regexp.MatchString(`^\d{12}$`, code); !match {
			logger.Warn(r.Context(), "msg", "invalid guide code")
			res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		logger.WithLogFieldsInRequest(r, "guide_code", code)

		gp, err := guide_process_provider.Get().GetGuideProcessByCode(r.Context(), code)
		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to fetch guide process")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}

		if gp.ID == 0 {
			logger.Info(r.Context(), "msg", "guide not found")
			res.Message = i18n.Get(r, i18n.MsgGuideNotFound)
			status := http.StatusNotFound
			response.WriteJSON(w, r, res, status)
			return
		}

		res.Data = gp
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type GuideInput struct {
	Code string `json:"code"`
}

type GuideOutput struct {
	ID int `json:"id"`
}

func CreateGuideProcess() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		logger.Info(r.Context(), "msg", "handler.GetGuideProcessByCode_start")
		defer logger.Info(r.Context(), "msg", "handler.GetGuideByCode_end")

		res := response.Response{}
		var input GuideInput
		// Decode JSON body into struct
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Warn(r.Context(), "msg", "missing guide code")
			res.Message = i18n.Get(r, i18n.MsgGuideRequired)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		// Optional: Validate fields
		if input.Code == "" {
			logger.Warn(r.Context(), "msg", "missing guide code")
			res.Message = i18n.Get(r, i18n.MsgGuideRequired)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		if match, _ := regexp.MatchString(`^\d{12}$`, input.Code); !match {
			logger.Warn(r.Context(), "msg", "invalid guide code")
			res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		guide, err := guide_provider.Get().GetGuide(r.Context(), input.Code)

		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to fetch guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}

		data := &GuideOutput{}
		res.Data = data

		if guide.ID == "" {
			logger.Info(r.Context(), "msg", "guide not found")
			res.Message = i18n.Get(r, i18n.MsgGuideNotFound)
			status := http.StatusNotFound
			response.WriteJSON(w, r, res, status)
			return
		}

		id, err := guide_process_provider.Get().CreateGuide(r.Context(), guide)

		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to create guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}
		data.ID = id
		return
	})
}
