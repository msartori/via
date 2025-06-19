package model

import (
	"time"
	"via/internal/ent"
)

type Guide struct {
	ID         int       `json:"id"`
	ViaGuideID string    `json:"viaGuideId"`
	Recipient  string    `json:"recipient"`
	OperatorID int       `json:"operatorId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func FromEntGuide(guide ent.Guide) Guide {
	return Guide{
		ID:         guide.ID,
		ViaGuideID: guide.ViaGuideID,
		Recipient:  guide.Recipient,
		Status:     guide.Status,
		CreatedAt:  guide.CreatedAt,
		UpdatedAt:  guide.UpdatedAt,
	}
}
