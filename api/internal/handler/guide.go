package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"via/internal/config"
	custom_error "via/internal/error"
	"via/internal/log"
	guide_provider "via/internal/provider/guide"
	util_response "via/internal/util/response"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
)

func GetGuide(biz config.Bussiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := util_response.Response{}

		logger := log.Get()

		id := chi.URLParam(r, "id")
		if id == "" {
			logger.Warn(ctx, "msg", "missing guide ID")
			response.Message = "El código de guía es requerido."
			util_response.ResponderJSON(w, r, response, http.StatusBadRequest)
			return
		}

		logger.Info(ctx, "msg", "starting guide lookup", "guideId", id)

		guide, err := guide_provider.Get().GetGuide(ctx, id)
		if err != nil {
			logger.Error(ctx, err, "msg", "failed to fetch guide", "guideId", id)
			var httpErr custom_error.HTTPError
			response.Message = "Error interno del sevidor."
			status := http.StatusInternalServerError
			if errors.As(err, &httpErr) {
				response.Message = httpErr.Message
				status = httpErr.StatusCode
			}
			util_response.ResponderJSON(w, r, response, status)
			return
		}
		data := &guide

		logger.Info(ctx, "msg", "guide found", "guideId", guide.ID)
		response.Data = data
		if guide.Destination.ID != biz.ViaBranch {
			response.Message = fmt.Sprintf("La guía solicitada corresponde a la sucursal %s", guide.Destination.Description)
			util_response.ResponderJSON(w, r, response, http.StatusOK)
			return
		}
		// in transit
		if slices.Contains(strings.Split(biz.PendingStatus, ","), guide.Status) {
			response.Message = "La guía solicitada se encuentra en tránsito, por favor vuelva más tarde."
			util_response.ResponderJSON(w, r, response, http.StatusOK)
			return
		}
		// delivered
		if guide.Status == biz.DeliveredStatus {
			response.Message = "La guía solicitada ya se ha entregado."
			util_response.ResponderJSON(w, r, response, http.StatusOK)
			return
		}
		// ready to withdraw
		if guide.Status == biz.WithdrawStatus {
			data.EnabledToWithdraw = true
			response.Message = "La guía solicitada esta disponible para su retiro."
		} else {
			// non processable status
			response.Message = "La guía solicitada no esta disponible."
		}
		util_response.ResponderJSON(w, r, response, http.StatusOK)
	})
}
