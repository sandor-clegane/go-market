package entities

import "time"

type Withdraw struct {
	Order       int       `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      string    `json:"user_id"`
}
