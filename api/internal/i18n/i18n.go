package i18n

import (
	"fmt"
	"net/http"
)

const (
	MsgGuideRequired           = "guide_required"
	MsgGuideInvalid            = "guide_invalid"
	MsgGuideNotFound           = "guide_not_found"
	MsgInternalServerError     = "internal_error"
	MsgOtherBranch             = "other_branch"
	MsgInTransit               = "in_transit"
	MsgDelivered               = "delivered"
	MsgWithdrawAvailable       = "withdraw_available"
	MsgNotAvailable            = "not_available"
	MsgRequestTimeout          = "request_timeout"
	MsgRequestCanceledByClient = "request_canceled_client"
	MsgUnexpectedContextError  = "unexpected_context_error"
)

var messages = map[string]map[string]string{
	"es": {
		MsgGuideRequired:           "El código de guía es requerido.",
		MsgGuideInvalid:            "El código de guía es inválido.",
		MsgGuideNotFound:           "Guía no econtrada.",
		MsgInternalServerError:     "Error interno del servidor.",
		MsgOtherBranch:             "La guía solicitada corresponde a la sucursal %s.",
		MsgInTransit:               "La guía solicitada se encuentra en tránsito, por favor vuelva más tarde.",
		MsgDelivered:               "La guía solicitada ya se ha entregado.",
		MsgWithdrawAvailable:       "La guía solicitada está disponible para su retiro.",
		MsgNotAvailable:            "La guía solicitada no está disponible.",
		MsgRequestTimeout:          "Tiempo de espera de solicitud agotado.",
		MsgRequestCanceledByClient: "Solicitud cancelada por el cliente.",
		MsgUnexpectedContextError:  "Error de contexto inesperado.",
	},
	"en": {
		MsgRequestTimeout:          "Request timeout.",
		MsgRequestCanceledByClient: "Request canceled by client.",
		MsgUnexpectedContextError:  "Unexpected context error.",
	},
}

func GetWithLang(lang, key string, args ...interface{}) string {
	if msg, ok := messages[lang][key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}
	// Fallback
	return key
}

func Get(r *http.Request, key string, args ...interface{}) string {
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
