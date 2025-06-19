package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"via/internal/i18n"
	"via/internal/log"
	response "via/internal/response"
)

func getJsonBody(w http.ResponseWriter, r *http.Request, input any) bool {
	res := response.Response[any]{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Get().Error(r.Context(), err, "msg", "failed to decode body")
		res.Message = i18n.Get(r, i18n.MsgInternalServerError)
		response.WriteJSON(w, r, res, http.StatusInternalServerError)
		return false
	}
	return true
}

func isValidViaGuideId(w http.ResponseWriter, r *http.Request, viaGuideId string) bool {
	res := response.Response[any]{}
	if viaGuideId == "" {
		log.Get().Warn(r.Context(), "msg", "missing via guide id")
		res.Message = i18n.Get(r, i18n.MsgGuideRequired)
		response.WriteJSON(w, r, res, http.StatusBadRequest)
		return false
	}

	if match, _ := regexp.MatchString(`^\d{12}$`, viaGuideId); !match {
		log.Get().Warn(r.Context(), "msg", "invalid guide code")
		res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
		response.WriteJSON(w, r, res, http.StatusBadRequest)
		return false
	}
	return true
}

func isGuideNotFound(w http.ResponseWriter, r *http.Request, id string) bool {
	res := response.Response[any]{}
	if id == "" {
		log.Get().Warn(r.Context(), "msg", "guide not found")
		res.Message = i18n.Get(r, i18n.MsgGuideNotFound)
		status := http.StatusNotFound
		response.WriteJSON(w, r, res, status)
		return false
	}
	return true
}

func isFailedToFetchGuide(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		res := response.Response[any]{}
		log.Get().Error(r.Context(), err, "msg", "failed to fetch guide")
		res.Message = i18n.Get(r, i18n.MsgInternalServerError)
		status := http.StatusInternalServerError
		response.WriteJSON(w, r, res, status)
		return true
	}
	return false
}
