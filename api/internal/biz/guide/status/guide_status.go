package biz_guide_status

import (
	"fmt"
	biz_config "via/internal/biz/config"
	biz_language "via/internal/biz/language"
)

// Status string constants
const (
	INITIAL                    = "initial"
	PENDING_RECIPIENT_IDENTIFY = "pendingRecipientIdentify"
	RECIPIENT_IDENTIFIED       = "recipientIdentified"
	PENDING_PAYMENT            = "pendingPayment"
	PAID                       = "paymentProcessed"
	PENDING_COUNTER_DELIVERY   = "pendingCounterDelivery"
	PENDING_WAREHOUSE_DELIVERY = "pendingWarehouseDelivery"
	PARTIAL_DELIVERED          = "partialDelivered"
	DELIVERED                  = "delivered"
	ON_HOLD                    = "onHold"
	SUSPENDED                  = "suspended"
	PREVIOUS                   = "previous"
)

// messages holds the localized descriptions
var messages = map[string]map[string]string{
	biz_language.ES: {
		INITIAL:                    "Inicial",
		PENDING_RECIPIENT_IDENTIFY: "Pendiente de identificación de destinatario",
		RECIPIENT_IDENTIFIED:       "Destinatario identificado",
		PENDING_PAYMENT:            "Pendiente de pago",
		PAID:                       "Pago realizado",
		PENDING_COUNTER_DELIVERY:   "Pendiente de retiro por mostrador",
		PENDING_WAREHOUSE_DELIVERY: "Pendiente de retiro por depósito",
		PARTIAL_DELIVERED:          "Parcialmente entregado",
		DELIVERED:                  "Entregado",
		ON_HOLD:                    "En espera",  // custormer can bring it back
		SUSPENDED:                  "Suspendida", // indeterminent, only operator can bring it back
	},
}

// Description returns the localized description of a status for a given language.
// If lang or status is unknown, it returns the status key itself.
func GetStatusDescription(lang, status string) string {
	if desc, ok := messages[lang][status]; ok {
		return desc
	}
	if lang != biz_language.DEFAULT {
		if desc, ok := messages[biz_language.DEFAULT][status]; ok {
			return desc
		}
	}

	return fmt.Sprintf("Unknown status: %s", status)
}

var nextStatus = map[string][]string{
	INITIAL:                    {ON_HOLD, SUSPENDED, PENDING_RECIPIENT_IDENTIFY},
	PENDING_RECIPIENT_IDENTIFY: {ON_HOLD, SUSPENDED, RECIPIENT_IDENTIFIED},
	//pending payment only if guide requires payment
	RECIPIENT_IDENTIFIED + biz_config.PAID_SHIPPING:       {ON_HOLD, SUSPENDED, PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY},
	RECIPIENT_IDENTIFIED + biz_config.PAID_ON_DESTINATION: {ON_HOLD, SUSPENDED, PENDING_PAYMENT},
	PENDING_PAYMENT:            {ON_HOLD, SUSPENDED, PAID},
	PAID:                       {ON_HOLD, SUSPENDED, PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY},
	PENDING_COUNTER_DELIVERY:   {ON_HOLD, SUSPENDED, PARTIAL_DELIVERED, DELIVERED},
	PENDING_WAREHOUSE_DELIVERY: {ON_HOLD, SUSPENDED, PARTIAL_DELIVERED, DELIVERED},
	PARTIAL_DELIVERED:          {},         //final status
	DELIVERED:                  {},         //final status
	ON_HOLD:                    {PREVIOUS}, //previous status (only for customer)
	SUSPENDED:                  {PREVIOUS}, //previous status (only for operator)
}

func GetNextStatus(currentStatus string, history []string, payment string) []string {
	if nStatus, found := nextStatus[currentStatus]; found {
		previous := ""
		if nStatus[0] == PREVIOUS {
			if len(history) > 1 {
				previous = history[len(history)-2]
			}
		}
		for i := len(history) - 3; i >= 0 && (previous == ON_HOLD || previous == SUSPENDED); i-- {
			previous = history[i]
		}
		if previous != "" {
			return []string{previous}
		}
		return nStatus
	} else {
		if nStatus, found = nextStatus[currentStatus+payment]; found {
			return nStatus
		}
		return []string{}
	}
}

func GetStatusState(status string) string {
	if status == SUSPENDED {
		return "error"
	}
	if status == ON_HOLD {
		return "warn"
	}
	return "ok"
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
		PENDING_COUNTER_DELIVERY == status ||
		PENDING_WAREHOUSE_DELIVERY == status
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
		PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY}
}

func GetOperatorStatus() []string {
	return []string{INITIAL, PENDING_RECIPIENT_IDENTIFY, RECIPIENT_IDENTIFIED, PENDING_PAYMENT, PAID,
		PENDING_COUNTER_DELIVERY, PENDING_WAREHOUSE_DELIVERY, ON_HOLD, SUSPENDED}
}
