package model

import "time"

type GuideHistory struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"-"`
}
