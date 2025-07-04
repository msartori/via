package i18n

import (
	"fmt"
	"net/http"
	"strings"
	"via/internal/response"
)

const (
	MsgGuideRequired             = "guide_required"
	MsgGuideInvalid              = "guide_invalid"
	MsgGuideNotFound             = "guide_not_found"
	MsgInternalServerError       = "internal_error"
	MsgRequestTimeout            = "request_timeout"
	MsgRequestCanceledByClient   = "request_canceled_client"
	MsgUnexpectedContextError    = "unexpected_context_error"
	MsgTooManyRequestsError      = "too_many_requests_error"
	MsgOperatorInvalid           = "operator_invalid"
	MsgOperatorUnauthorized      = "operator_unauthorized"
	MsgAuthStateNotFound         = "auth_state_not_found"
	MsgAuthFailedToExchangeToken = "auth_failed_to_exchange_token"
	MsgAuthFailedToGetUserInfo   = "auth_failed_to_get_user_info"
)

var messages = map[string]map[string]string{
	"es": {
		MsgGuideRequired:             "El código de guía es requerido.",
		MsgGuideInvalid:              "El código de guía es inválido.",
		MsgGuideNotFound:             "Guía no econtrada.",
		MsgInternalServerError:       "Error interno del servidor.",
		MsgRequestTimeout:            "Tiempo de espera de solicitud agotado.",
		MsgRequestCanceledByClient:   "Solicitud cancelada por el cliente.",
		MsgUnexpectedContextError:    "Error de contexto inesperado.",
		MsgTooManyRequestsError:      "Demasiadas solicitudes, por favor intente más tarde.",
		MsgOperatorInvalid:           "El Id de Operador es inválido",
		MsgOperatorUnauthorized:      "Operador no autorizado.",
		MsgAuthStateNotFound:         "Estado de autenticación no encontrado.",
		MsgAuthFailedToExchangeToken: "Error al intercambiar el token de autenticación.",
		MsgAuthFailedToGetUserInfo:   "Error al obtener la información del usuario.",
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

func Get(r *http.Request, key string, args ...any) string {
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
