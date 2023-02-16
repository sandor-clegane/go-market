package utils

import (
	"strconv"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities/customErrors"

	"github.com/theplant/luhn"
)

func ValidateNumber(order string) (int, error) {
	orderNumber, err := strconv.Atoi(order)
	if err != nil {
		return 0, err
	}
	if !luhn.Valid(orderNumber) {
		return 0, customErrors.NewInvalidOrderNumberFormatError(orderNumber)
	}
	return orderNumber, nil
}
