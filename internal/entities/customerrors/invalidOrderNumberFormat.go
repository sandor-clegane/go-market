package customerrors

import "fmt"

type InvalidOrderNumberFormatError struct {
	Number int
}

func (iof *InvalidOrderNumberFormatError) Error() string {
	return fmt.Sprintf("order with number %d has invalid format", iof.Number)
}

func NewInvalidOrderNumberFormatError(number int) error {
	return &InvalidOrderNumberFormatError{
		Number: number,
	}
}
