package biz_guide_status

import "fmt"

// Status string constants
const (
	INITIAL                    = "initial"
	PENDING_RECIPIENT_IDENTIFY = "pendingRecipientIdentify"
	RECIPIENT_IDENTIFIED       = "recipientIdentified"
	PENDING_PAYMENT            = "pendingPayment"
	PAID                       = "paymentProcessed"
	//WITHDRAW_ENABLED           = "withdrawEnabled"
	PENDING_COUNTER_DELIVERY   = "pendingCounterDelivery"
	PENDING_WAREHOUSE_DELIVERY = "pendingWarehouseDelivery"
	PARTIAL_DELIVERED          = "partialDelivered"
	DELIVERED                  = "delivered"
	ON_HOLD                    = "onHold"
	SUSPENDED                  = "suspended"
	RECOVERED                  = "recovered"
)

// messages holds the localized descriptions
var messages = map[string]map[string]string{
	"es": {
		INITIAL:                    "Inicial",
		PENDING_RECIPIENT_IDENTIFY: "Pendiente de identificación de destinatario",
		RECIPIENT_IDENTIFIED:       "Destinatario identificado",
		PENDING_PAYMENT:            "Pendiente de pago",
		PAID:                       "Pago realizado",
		//WITHDRAW_ENABLED:           "Habilitado para retiro",
		PENDING_COUNTER_DELIVERY:   "Pendiente de retiro por mostrador",
		PENDING_WAREHOUSE_DELIVERY: "Pendiente de retiro por depósito",
		PARTIAL_DELIVERED:          "Parcialmente entregado",
		DELIVERED:                  "Entregado",
		ON_HOLD:                    "En espera",  // custormer can bring it back
		SUSPENDED:                  "Suspendida", // indeterminent, only operator can bring it back
		RECOVERED:                  "Recuperada", // was ON-HOLD or SUSPENDED
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
		PENDING_RECIPIENT_IDENTIFY == status ||
		RECIPIENT_IDENTIFIED == status ||
		PENDING_PAYMENT == status ||
		PAID == status ||
		//WITHDRAW_ENABLED == status ||
		PENDING_COUNTER_DELIVERY == status ||
		PENDING_WAREHOUSE_DELIVERY == status ||
		RECOVERED == status
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

func GetMonitorStatus() []string {
	return []string{INITIAL, PENDING_RECIPIENT_IDENTIFY, RECIPIENT_IDENTIFIED, PENDING_PAYMENT, PAID,
		PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY, RECOVERED}
}

func GetOperatorStatus() []string {
	return []string{INITIAL, PENDING_RECIPIENT_IDENTIFY, RECIPIENT_IDENTIFIED, PENDING_PAYMENT, PAID,
		PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY, RECOVERED}
}
