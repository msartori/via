package model

import (
	biz_guide_status "via/internal/biz/guide/status"
)

type MonitorEvent struct {
	GuideId   string `json:"guideId"`
	Recipient string `json:"recipient"`
	Status    string `json:"status"`
	Highlight bool   `json:"highlight"`
}

const (
	MSG_WAIT      = "msgWait"
	MSG_COUNTER   = "msgCounter"
	MSG_WAREHOUSE = "msgWarehouse"
)

var monitorEventStatus = map[string]map[string]string{
	"es": {
		MSG_WAIT:      "Aguarde por favor",
		MSG_COUNTER:   "Presentarse en mostrador",
		MSG_WAREHOUSE: "Presentarse en depósito",
	},
}

var monitorEventStatusByGuideStatus = map[string]string{
	biz_guide_status.INITIAL:                    MSG_WAIT,
	biz_guide_status.PENDING_RECIPIENT_IDENTIFY: MSG_COUNTER,
	biz_guide_status.RECIPIENT_IDENTIFIED:       MSG_WAIT,
	biz_guide_status.PENDING_PAYMENT:            MSG_COUNTER,
	biz_guide_status.PAID:                       MSG_WAIT,
	biz_guide_status.PENDING_COUNTER_DELIVERY:   MSG_COUNTER,
	biz_guide_status.RECOVERED:                  MSG_WAIT,
}

func GetMonitorEventStatusDescriptionByGuideStatus(lang, guideStatus string) string {
	status, found := monitorEventStatusByGuideStatus[guideStatus]
	if !found {
		return MSG_WAIT
	}
	return monitorEventStatus[lang][status]
}

func GetHighlightByGuideStatus(guideStatus string) bool {
	return monitorEventStatusByGuideStatus[guideStatus] != MSG_WAIT
}
