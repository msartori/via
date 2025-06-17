package handler

import (
	"net/http"
	"regexp"
	"strings"
	"via/internal/config"
	"via/internal/i18n"
	"via/internal/log"
	via_guide_provider "via/internal/provider/via/guide"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
)

func GetGuide(biz config.Bussiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		logger.Info(r.Context(), "msg", "handler.GetGuide_start")
		defer logger.Info(r.Context(), "msg", "handler.GetGuide_end")

		res := response.Response{}
		id := chi.URLParam(r, "id")

		if id == "" {
			logger.Warn(r.Context(), "msg", "missing guide ID")
			res.Message = i18n.Get(r, i18n.MsgGuideRequired)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		if match, _ := regexp.MatchString(`^\d{12}$`, id); !match {
			logger.Warn(r.Context(), "msg", "invalid guide ID")
			res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		logger.WithLogFieldsInRequest(r, "guideId", id)

		guide, err := via_guide_provider.Get().GetGuide(r.Context(), id)
		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to fetch guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}

		data := &guide
		res.Data = data

		if guide.ID == "" {
			logger.Info(r.Context(), "msg", "guide not found")
			res.Message = i18n.Get(r, i18n.MsgGuideNotFound)
			status := http.StatusNotFound
			response.WriteJSON(w, r, res, status)
			return
		}

		logger.Info(r.Context(), "msg", "guide found")

		if guide.Destination.ID != biz.ViaBranch {
			res.Message = i18n.Get(r, i18n.MsgOtherBranch, guide.Destination.Description)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		if slices.Contains(strings.Split(biz.PendingStatus, ","), guide.Status) {
			res.Message = i18n.Get(r, i18n.MsgInTransit)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		if guide.Status == biz.DeliveredStatus {
			res.Message = i18n.Get(r, i18n.MsgDelivered)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		if guide.Status == biz.WithdrawStatus {
			data.EnabledToWithdraw = true
			res.Message = i18n.Get(r, i18n.MsgWithdrawAvailable)
		} else {
			res.Message = i18n.Get(r, i18n.MsgNotAvailable)
		}

		response.WriteJSON(w, r, res, http.StatusOK)
	})
}
