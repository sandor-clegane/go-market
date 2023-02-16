package customErrors

import "fmt"

type ExistedOrderError struct {
	Number int
	UserID string
}

func (eo *ExistedOrderError) Error() string {
	return fmt.Sprintf("order with number %d was saved by user with ID %s", eo.Number, eo.UserID)
}

func NewExistedOrderError(number int, userID string) error {
	return &ExistedOrderError{
		Number: number,
		UserID: userID,
	}
}
