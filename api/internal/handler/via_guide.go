package handler

import (
	"net/http"
	"slices"
	"strings"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/config"
	"via/internal/log"
	guide_provider "via/internal/provider/guide"
	via_guide_provider "via/internal/provider/via/guide"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

type GetGuideToWithdrawOutput struct {
	EnabledToWithdraw bool   `json:"enabledToWithdraw"`
	WithdrawMessage   string `json:"withdrawMessage"`
}

func GetGuideToWithdraw(biz config.Bussiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[*GetGuideToWithdrawOutput]{}
		viaGuideId := chi.URLParam(r, "viaGuideId")

		if ok := isValidViaGuideId(w, r, viaGuideId); !ok {
			return
		}

		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "viaGuideId", viaGuideId)

		viaGuide, err := via_guide_provider.Get().GetGuide(r.Context(), viaGuideId)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}
		data := &GetGuideToWithdrawOutput{EnabledToWithdraw: false}
		res.Data = data
		if viaGuide.ID == "" {
			data.WithdrawMessage = getWithDrawMessage(r, notFound, viaGuideId)
			logger.Info(r.Context(), "msg", "guide not found")
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if viaGuide.Status == biz.DeliveredStatus {
			logger.Info(r.Context(), "msg", "guide already delivered")
			data.WithdrawMessage = getWithDrawMessage(r, delivered, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if viaGuide.Destination.ID != biz.ViaBranch {
			logger.Info(r.Context(), "msg", "guide is not in this branch")
			data.WithdrawMessage = getWithDrawMessage(r, wrongBranch, viaGuideId, viaGuide.Destination.Description)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if slices.Contains(strings.Split(biz.PendingStatus, ","), viaGuide.Status) {
			logger.Info(r.Context(), "msg", "guide is in via pending status")
			data.WithdrawMessage = getWithDrawMessage(r, pending, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if viaGuide.Status == biz.WithdrawStatus {
			logger.Info(r.Context(), "msg", "validating withdraw status")
			guide, err := guide_provider.Get().GetGuideByViaGuideId(r.Context(), viaGuideId)
			if isFailedToFetchGuide(w, r, err) {
				return
			}
			//guide found
			if guide.ID != 0 {
				if biz_guide_status.IsDelivered(guide.Status) {
					logger.Info(r.Context(), "msg", "guide already delivered, not impacted in via system")
					data.WithdrawMessage = getWithDrawMessage(r, delivered, viaGuideId)
					response.WriteJSON(w, r, res, http.StatusOK)
					return
				}
				if biz_guide_status.IsInProcess(guide.Status) {
					logger.Info(r.Context(), "msg", "guide in process")
					data.WithdrawMessage = getWithDrawMessage(r, alreadyInProcess, viaGuideId)
					response.WriteJSON(w, r, res, http.StatusOK)
					return
				}
				if biz_guide_status.IsEnabledToWithdraw(guide.Status) {
					logger.Info(r.Context(), "msg", "guide enabled to start withdraw process")
					data.EnabledToWithdraw = true
					data.WithdrawMessage = getWithDrawMessage(r, enabledToWithdraw, viaGuideId)
					response.WriteJSON(w, r, res, http.StatusOK)
					return
				}
			}
			logger.Info(r.Context(), "msg", "guide enabled to start withdraw process")
			data.EnabledToWithdraw = true
			data.WithdrawMessage = getWithDrawMessage(r, enabledToWithdraw, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		logger.Info(r.Context(), "msg", "guide not available to start withdraw process", "guide", viaGuide)
		data.WithdrawMessage = getWithDrawMessage(r, notAvailable, viaGuideId)
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}
