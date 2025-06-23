package model

import (
	"time"
)

type Guide struct {
	ID         int       `json:"id"`
	ViaGuideID string    `json:"viaGuideId"`
	Recipient  string    `json:"recipient"`
	Operator   Operator  `json:"operator"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
