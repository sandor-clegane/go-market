package withdrawservice

import (
	"context"
	"sort"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customerrors"
	"github.com/sandor-clegane/go-market/internal/entities/mappers"
	"github.com/sandor-clegane/go-market/internal/storage/orderstorage"
	"github.com/sandor-clegane/go-market/internal/storage/withdrawstorage"
	"github.com/sandor-clegane/go-market/internal/utils"
)

type withdrawServiceImpl struct {
	withdrawRepository withdrawstorage.WithdrawStorage
	orderRepository    orderstorage.OrderStorage
}

func New(withdrawRepository withdrawstorage.WithdrawStorage, orderRepository orderstorage.OrderStorage) WithdrawService {
	return &withdrawServiceImpl{
		withdrawRepository: withdrawRepository,
		orderRepository:    orderRepository,
	}
}

func (w *withdrawServiceImpl) CreateWithdraw(ctx context.Context, withdrawRequest entities.WithdrawRequest, userID string) error {
	order, err := utils.ValidateNumber(withdrawRequest.Order)
	if err != nil {
		return err
	}
	totalUserAccrual, err := w.orderRepository.GetTotalAccrualAmountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if totalUserAccrual < withdrawRequest.Sum {
		return customerrors.NewLimitExceededError(withdrawRequest.Sum, totalUserAccrual, userID)
	}
	return w.withdrawRepository.InsertWithdraw(ctx, mappers.MapToWithdraw(withdrawRequest, order, userID))
}

func (w *withdrawServiceImpl) GetBalanceInfoByID(ctx context.Context, userID string) (entities.BalanceRequest, error) {
	totalUserAccrual, err := w.orderRepository.GetTotalAccrualAmountByUserID(ctx, userID)
	if err != nil {
		return entities.BalanceRequest{}, err
	}
	totalUserWithdraw, err := w.withdrawRepository.GetTotalWithdrawnByUserID(ctx, userID)
	if err != nil {
		return entities.BalanceRequest{}, err
	}
	return mappers.MapToBalanceRequest(totalUserAccrual, totalUserWithdraw), nil
}

func (w *withdrawServiceImpl) GetWithdrawsInfoByID(ctx context.Context, userID string) ([]entities.WithdrawDTO, error) {
	withdrawList, err := w.withdrawRepository.GetAllWithdrawsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	sort.Slice(withdrawList, func(i, j int) bool {
		return withdrawList[j].ProcessedAt.Before(withdrawList[i].ProcessedAt)
	})
	return mappers.MapWithdrawListToWithdrawDTOList(withdrawList), nil
}
