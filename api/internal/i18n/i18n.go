package i18n

import (
	"fmt"
	"net/http"
	"strings"
	"via/internal/response"
)

const (
	MsgGuideRequired           = "guide_required"
	MsgGuideInvalid            = "guide_invalid"
	MsgGuideNotFound           = "guide_not_found"
	MsgInternalServerError     = "internal_error"
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
			if strings.Contains(msg, "%") {
				return fmt.Sprintf(msg, args...)
			}
		}
		return msg
	}
	// Fallback
	return key
}

func Get(r *http.Request, key string, args ...interface{}) string {
	lang := response.GetLanguage(r)
	if msg, ok := messages[lang][key]; ok {
		if len(args) > 0 {
			if strings.Contains(msg, "%") {
				return fmt.Sprintf(msg, args...)
			}
		}
		return msg
	}
	// Fallback
	return key
}
