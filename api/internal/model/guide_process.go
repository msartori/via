package model

import "time"

type GuideProcess struct {
	ID         int       `json:"id"`
	Code       string    `json:"code"`
	Recipient  string    `json:"recipient"`
	OperatorID int       `json:"operatorId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
