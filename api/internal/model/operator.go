package model

type Operator struct {
	ID      int    `json:"id"`
	Account string `json:"account"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
