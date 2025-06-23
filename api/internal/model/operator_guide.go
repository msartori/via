package model

type OperatorGuide struct {
	GuideId   int      `json:"guideId"`
	Recipient string   `json:"recipient"`
	Status    string   `json:"status"`
	Operator  Operator `json:"operator"`
}
