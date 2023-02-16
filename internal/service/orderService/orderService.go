package orderService

import (
	"context"
	"sort"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
	"Gophermarket/go-musthave-diploma-tpl/internal/entities/mappers"
	"Gophermarket/go-musthave-diploma-tpl/internal/storage/orderStorage"
	"Gophermarket/go-musthave-diploma-tpl/internal/utils"
)

type orderServiceImpl struct {
	orderRepository orderStorage.OrderStorage
}

func New(orderRepository orderStorage.OrderStorage) OrderService {
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
