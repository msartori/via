package biz_config

import (
	"fmt"
	biz_language "via/internal/biz/language"
)

type Bussiness struct {
	ViaBranch       string `env:"VIA_BRANCH" envDefault:"123" json:"viaBranch"`
	WithdrawStatus  string `env:"WITHDRAW_STATUS" envDefault:"CRR" json:"withdrawStatus"`
	DeliveredStatus string `env:"DELIVERED_STATUS" envDefault:"ENT" json:"deliveredStatus"`
	PendingStatus   string `env:"PENDING_STATUS" envDefault:"ASP,CPO,ORI,PTE" json:"pendingStatus"`
	HomeDelivery    string `env:"HOME_DELIVERY" envDefault:"CD06" json:"homeDelivery"`
}

const (
	PAID_SHIPPING       = "P"
	PAID_ON_DESTINATION = "D"
)

var paymentDescription = map[string]map[string]string{
	biz_language.ES: {
		PAID_SHIPPING:       "Origen",
		PAID_ON_DESTINATION: "Destino",
	}}

func GetPaymentDescription(lang, payment string) string {
	if desc, ok := paymentDescription[lang][payment]; ok {
		return desc
	}
	if lang != biz_language.DEFAULT {
		if desc, ok := paymentDescription[biz_language.DEFAULT][payment]; ok {
			return desc
		}
	}

	return fmt.Sprintf("Unknown payment: %s", payment)
}
