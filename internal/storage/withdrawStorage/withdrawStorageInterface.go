package withdrawStorage

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var ws WithdrawStorage = &withdrawStorageImpl{}

type WithdrawStorage interface {
	InsertWithdraw(ctx context.Context, withdraw entities.Withdraw) error
	GetTotalWithdrawnByUserID(ctx context.Context, userID string) (float32, error)
	GetAllWithdrawsByUserID(ctx context.Context, userID string) ([]entities.Withdraw, error)
}
