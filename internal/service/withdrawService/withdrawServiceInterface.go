package withdrawService

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var _ WithdrawService = &withdrawServiceImpl{}

type WithdrawService interface {
	CreateWithdraw(ctx context.Context, withdrawRequest entities.WithdrawRequest, userID string) error
	GetBalanceInfoByID(ctx context.Context, userID string) (entities.BalanceRequest, error)
	GetWithdrawsInfoByID(ctx context.Context, userID string) ([]entities.WithdrawDTO, error)
}