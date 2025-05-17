package model_guide

import model_destination "via/internal/model/destination"

type Guide struct {
	ID          int                           `json:"id"`
	Destination model_destination.Destination `json:"destination"`
}
