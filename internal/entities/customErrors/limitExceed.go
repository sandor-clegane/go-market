package customErrors

import "fmt"

type LimitExceededError struct {
	Sum          float32
	TotalAccrual float32
	UserID       string
}

func (le *LimitExceededError) Error() string {
	return fmt.Sprintf("sum %f exceededs limit, for user %s accessible limit is %f", le.Sum, le.UserID, le.TotalAccrual)
}

func NewLimitExceededError(sum float32, totalAccrual float32, userID string) error {
	return &LimitExceededError{
		Sum:          sum,
		TotalAccrual: totalAccrual,
		UserID:       userID,
	}
}
