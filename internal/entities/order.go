package entities

import (
	"fmt"
	"time"
)

type Order struct {
	Number     int         `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
	UserID     string      `json:"user_id"`
}

type OrderStatus uint8

const (
	NEW OrderStatus = iota + 1
	PROCESSING
	INVALID
	PROCESSED
)

func (o OrderStatus) String() string {
	switch o {
	case NEW:
		return "NEW"
	case PROCESSING:
		return "PROCESSING"
	case INVALID:
		return "INVALID"
	case PROCESSED:
		return "PROCESSED"
	default:
		return fmt.Sprintf("%d", int(o))
	}
}
