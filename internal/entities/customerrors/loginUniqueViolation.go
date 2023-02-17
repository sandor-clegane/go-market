package customerrors

import "fmt"

type UserLoginUniqueViolation struct {
	Err   error
	Login string
}

func (ve *UserLoginUniqueViolation) Error() string {
	return fmt.Sprintf("user with login %s already exists", ve.Login)
}

func (ve *UserLoginUniqueViolation) Unwrap() error {
	return ve.Err
}

func NewUserLoginUniqueViolationError(login string, err error) error {
	return &UserLoginUniqueViolation{
		Login: login,
		Err:   err,
	}
}
