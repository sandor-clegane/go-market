package utils

import (
	"strconv"

	"github.com/sandor-clegane/go-market/internal/entities/customerrors"

	"github.com/theplant/luhn"
)

func ValidateNumber(order string) (int, error) {
	orderNumber, err := strconv.Atoi(order)
	if err != nil {
		return 0, err
	}
	if !luhn.Valid(orderNumber) {
		return 0, customerrors.NewInvalidOrderNumberFormatError(orderNumber)
	}
	return orderNumber, nil
}
