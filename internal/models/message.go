package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	Processed bool      `json:"processed"`
}
