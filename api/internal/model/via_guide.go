package model

type ViaGuide struct {
	ID                string         `json:"id"`
	Reference         string         `json:"reference"`
	Status            string         `json:"status"`
	Packages          int            `json:"packages"`
	Weight            float64        `json:"weight"`
	Payment           string         `json:"payment"`
	Route             string         `json:"route"`
	Date              string         `json:"date"`
	Sender            string         `json:"sender"`
	Recipient         string         `json:"recipient"`
	Destination       ViaDestination `json:"destination"`
	EnabledToWithdraw bool           `json:"enabledToWithdraw"`
}
