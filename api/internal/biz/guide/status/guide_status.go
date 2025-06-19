package biz_guide_status

import "fmt"

// Status string constants
const (
	INITIAL              = "initial"
	RECIPIENT_IDENTIFIED = "recipientIdentified"
	PAYMENT_PROCESSED    = "paymentProcessed"
	WITHDRAW_ENABLED     = "withdrawEnabled"
	COUNTER_DELIVERY     = "counterDelivery"
	WAREHOUSE_DELIVERY   = "warehouseDelivery"
	PARTIAL_DELIVERED    = "partialDelivered"
	DELIVERED            = "delivered"
	ON_HOLD              = "onHold"
)

// messages holds the localized descriptions
var messages = map[string]map[string]string{
	"es": {
		INITIAL:              "Inicial",
		RECIPIENT_IDENTIFIED: "Destinatario identificado",
		PAYMENT_PROCESSED:    "Pago realizado",
		WITHDRAW_ENABLED:     "Habilitado para retiro",
		COUNTER_DELIVERY:     "Retiro por mostrador",
		WAREHOUSE_DELIVERY:   "Retiro por depósito",
		PARTIAL_DELIVERED:    "Parcialmente entregado",
		DELIVERED:            "Entregado",
		ON_HOLD:              "En espera",
	},
}

// Description returns the localized description of a status for a given language.
// If lang or status is unknown, it returns the status key itself.
func Description(lang, status string) string {
	if desc, ok := messages[lang][status]; ok {
		return desc
	}
	return fmt.Sprintf("Unknown status: %s", status)
}

func IsEnabledToWithdraw(status string) bool {
	return ON_HOLD == status || PARTIAL_DELIVERED == status
}

func IsInProcess(status string) bool {
	return INITIAL == status ||
		RECIPIENT_IDENTIFIED == status ||
		PAYMENT_PROCESSED == status ||
		WITHDRAW_ENABLED == status ||
		COUNTER_DELIVERY == status ||
		WAREHOUSE_DELIVERY == status
}

func IsDelivered(status string) bool {
	return DELIVERED == status
}

func IsAbleToReInit(status string) bool {
	return ON_HOLD == status
}

func IsValidToCreateForWithdraw(status string) bool {
	return PARTIAL_DELIVERED == status
}
