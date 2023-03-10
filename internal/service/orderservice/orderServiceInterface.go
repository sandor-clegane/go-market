package orderservice

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var _ OrderService = (*orderServiceImpl)(nil)

type OrderService interface {
	CreateOrder(ctx context.Context, order, userID string) error
	GetAllOrdersByUserID(ctx context.Context, userID string) ([]entities.OrderDTO, error)
}
