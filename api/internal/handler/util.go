package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"via/internal/config"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/model"
	response "via/internal/response"
)

const (
	notFound          = "not_found"
	wrongBranch       = "wrong_branch"
	pending           = "pending"
	enabledToWithdraw = "enabled_to_withdraw"
	notAvailable      = "not_available"
	delivered         = "delivered"
	alreadyInProcess  = "already_in_process"
	inProcess         = "in_process"
	notAbleToProcess  = "not_able_to_process"
)

var messages = map[string]map[string]string{
	"es": {
		notFound:          "El código de guía ingresado [%s] no se encuentra, por favor verifique que sea correcto.",
		wrongBranch:       "El código de guía ingresado [%s] es correcto, pero el destino del envío no corresponde a esta sucursal. Por favor dirijase a la sucursal %s.",
		pending:           "El código de guía ingresado [%s] es correcto, pero el envío no ha llegado a destino. Por favor regrese en otro momento.",
		enabledToWithdraw: "El código de guía ingresado [%s] es correcto y el envío se encuentra disponible para ser retirado. Si Ud. es el destinatario deberá presentar su DNI para realizar el retiro. Presione Aceptar y aguarde a ser llamado por mostrador para verificar su identidad y proceder al retiro.",
		notAvailable:      "El código de guía ingresado [%s] es correcto, pero el envío no se encuentra disponible. Por favor consulte por mostrador.",
		delivered:         "El código de guía ingresado [%s] es correcto, pero el envío ya ha sido retirado.",
		alreadyInProcess:  "El código de guía ingresado [%s] es correcto, el envío esta en proceso de entrega, aguarde y será atendido.",
		inProcess:         "La guía [%s] está en proceso. Por favor aguarde a ser atendido.",
		notAbleToProcess:  "No es posible processar el retiro de la guía [%s], vuelva a consultar más tarde.",
	},
}

func getWithDrawMessage(r *http.Request, key string, args ...interface{}) string {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "es"
	}
	if msg, ok := messages[lang][key]; ok {
		if len(args) > 0 && strings.Contains(msg, "%") {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}
	// Fallback
	return key
}

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

func isInvalidViaGuideToWithdraw(viaGuide model.ViaGuide, biz config.Bussiness) bool {
	if viaGuide.Status == biz.DeliveredStatus {
		return false
	}
	if viaGuide.Destination.ID != biz.ViaBranch {
		return false
	}
	if slices.Contains(strings.Split(biz.PendingStatus, ","), viaGuide.Status) {
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
