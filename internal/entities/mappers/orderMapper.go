package mappers

import (
	"strconv"
	"time"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
)

func MapToOrder(number int, userID string) entities.Order {
	return entities.Order{
		Number:     number,
		UploadedAt: time.Now(),
		Status:     entities.NEW,
		UserID:     userID,
	}
}

func MapOrderListToOrderDTOList(orderList []entities.Order) []entities.OrderDTO {
	orderDTOList := make([]entities.OrderDTO, 0, len(orderList))
	var dto entities.OrderDTO
	for _, o := range orderList {
		dto = entities.OrderDTO{
			Number:     strconv.Itoa(o.Number),
			Status:     o.Status.String(),
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		orderDTOList = append(orderDTOList, dto)
	}
	return orderDTOList
}

func MapOrderResponseToOrder(orderResponse entities.OrderResponse) (entities.Order, error) {
	number, err := strconv.Atoi(orderResponse.Order)
	if err != nil {
		return entities.Order{}, err

	}
	order := entities.Order{
		Number: number,
	}
	if orderResponse.Status == "REGISTERED" || orderResponse.Status == "PROCESSING" {
		order.Status = entities.OrderStatus(2)
	} else if orderResponse.Status == "INVALID" {
		order.Status = entities.OrderStatus(3)
	} else {
		order.Status = entities.OrderStatus(4)
		order.Accrual = orderResponse.Accrual
	}
	return order, nil
}
