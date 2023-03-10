package orderservice

import (
	"context"
	"sort"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/mappers"
	"github.com/sandor-clegane/go-market/internal/storage/orderstorage"
	"github.com/sandor-clegane/go-market/internal/utils"
)

type orderServiceImpl struct {
	orderRepository orderstorage.OrderStorage
}

func New(orderRepository orderstorage.OrderStorage) OrderService {
	return &orderServiceImpl{
		orderRepository,
	}
}

func (o *orderServiceImpl) CreateOrder(ctx context.Context, order string, userID string) error {
	number, err := utils.ValidateNumber(order)
	if err != nil {
		return err
	}
	return o.orderRepository.InsertOrder(ctx, mappers.MapToOrder(number, userID))
}

func (o *orderServiceImpl) GetAllOrdersByUserID(ctx context.Context, userID string) ([]entities.OrderDTO, error) {
	orderList, err := o.orderRepository.GetAllOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	sort.Slice(orderList, func(i, j int) bool {
		return orderList[j].UploadedAt.Before(orderList[i].UploadedAt)
	})
	return mappers.MapOrderListToOrderDTOList(orderList), nil
}
