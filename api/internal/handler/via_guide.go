package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"via/internal/config"
	"via/internal/log"
	via_guide_provider "via/internal/provider/via/guide"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

const (
	notFound          = "not_found"
	wrongBranch       = "wrong_branch"
	pending           = "pending"
	enabledToWithdraw = "enabled_to_withdraw"
	notAvailable      = "not_available"
	delivered         = "delivered"
)

var messages = map[string]map[string]string{
	"es": {
		notFound:          "El código de guía ingresado [%s] no se encuentra, por favor verifique que sea correcto.",
		wrongBranch:       "El código de guía ingresado [%s] es correcto, pero el destino del envío no corresponde a esta sucursal. Por favor dirijase a la sucursal %s.",
		pending:           "El código de guía ingresado [%s] es correcto, pero el envío no ha llegado a destino. Por favor regrese en otro momento.",
		enabledToWithdraw: "El código de guía ingresado [%s] es correcto y el envío se encuentra disponible para ser retirado. Si Ud. es el destinatario deberá presentar su DNI para realizar el retiro. Presione Aceptar y aguarde a ser llamado por mostrador para verificar su identidad y proceder al retiro.",
		notAvailable:      "El código de guía ingresado [%s] es correcto, pero el envío no se encuentra disponible. Por favor consulte por mostrador.",
		delivered:         "El código de guía ingresado [%s] es correcto, pero el envío ya ha sido retirado.",
	},
}

type GetGuideToWithdrawOutput struct {
	EnabledToWithdraw bool   `json:"enabledToWithdraw"`
	WithdrawMessage   string `json:"withdrawMessage"`
}

func getWithDrawMessage(r *http.Request, key string, args ...interface{}) string {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "es"
	}
	if msg, ok := messages[lang][key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}
	// Fallback
	return key
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
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		if viaGuide.Status == biz.DeliveredStatus {
			data.WithdrawMessage = getWithDrawMessage(r, delivered, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if viaGuide.Destination.ID != biz.ViaBranch {
			data.WithdrawMessage = getWithDrawMessage(r, wrongBranch, viaGuideId, viaGuide.Destination.Description)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}
		if slices.Contains(strings.Split(biz.PendingStatus, ","), viaGuide.Status) {
			data.WithdrawMessage = getWithDrawMessage(r, pending, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		if viaGuide.Status == biz.WithdrawStatus {
			data.EnabledToWithdraw = true
			data.WithdrawMessage = getWithDrawMessage(r, enabledToWithdraw, viaGuideId)
			response.WriteJSON(w, r, res, http.StatusOK)
			return
		}

		data.WithdrawMessage = getWithDrawMessage(r, notAvailable, viaGuideId)
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}
