package model

import "time"

type OperatorGuide struct {
	GuideId    int       `json:"guideId"`
	ViaGuideId string    `json:"viaGuideId"`
	Recipient  string    `json:"recipient"`
	Status     string    `json:"status"`
	LastChange time.Time `json:"lastChange"`
	Payment    string    `json:"payment"`
	Operator   Operator  `json:"operator"`
	Selectable bool      `json:"selectable"`
}
