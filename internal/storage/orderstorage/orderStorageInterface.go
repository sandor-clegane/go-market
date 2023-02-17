package orderstorage

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var os OrderStorage = &orderStorageImpl{}

type OrderStorage interface {
	FindByNumber(ctx context.Context, number int) (entities.Order, error)
	GetAllOrdersByUserID(ctx context.Context, userID string) ([]entities.Order, error)
	GetTotalAccrualAmountByUserID(ctx context.Context, userID string) (float32, error)
	InsertOrder(ctx context.Context, order entities.Order) error
	StopSchedulerAndWorkerPool()
}
