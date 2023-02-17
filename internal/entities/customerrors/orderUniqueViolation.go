package customerrors

import (
	"fmt"

	"github.com/sandor-clegane/go-market/internal/entities"
)

type OrderViolationError struct {
	Err   error
	Order entities.Order
}

func (oe *OrderViolationError) Error() string {
	return fmt.Sprintf("order with number %d already saved", oe.Order.Number)
}

func (oe *OrderViolationError) Unwrap() error {
	return oe.Err
}

func NewOrderViolationError(order entities.Order, err error) error {
	return &OrderViolationError{
		Order: order,
		Err:   err,
	}
}
