package orderService

import (
	"context"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
)

var _ OrderService = &orderServiceImpl{}

type OrderService interface {
	CreateOrder(ctx context.Context, order, userID string) error
	GetAllOrdersByUserID(ctx context.Context, userID string) ([]entities.OrderDTO, error)
}
