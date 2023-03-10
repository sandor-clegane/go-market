package customerrors

import "fmt"

type LoginUniqueViolationError struct {
	Err   error
	Login string
}

func (ve *LoginUniqueViolationError) Error() string {
	return fmt.Sprintf("user with login %s already exists", ve.Login)
}

func (ve *LoginUniqueViolationError) Unwrap() error {
	return ve.Err
}

func NewUserLoginUniqueViolationError(login string, err error) error {
	return &LoginUniqueViolationError{
		Login: login,
		Err:   err,
	}
}
