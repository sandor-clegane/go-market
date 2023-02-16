package accrualService

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
)

const (
	getEndpoint = "api/orders"
)

type AccrualService struct {
	accrualSystemAddress string
}

func New(AccrualServiceAddress string) *AccrualService {
	return &AccrualService{
		accrualSystemAddress: AccrualServiceAddress,
	}
}

//GetOrderInfo - получение информации о расчёте начислений баллов лояльности.
func (ac *AccrualService) GetOrderInfo(orderID int) (orderInfo *entities.OrderResponse, err error) {
	targetURL := fmt.Sprintf("%s/%s/%d", ac.accrualSystemAddress, getEndpoint, orderID)

	res, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(orderInfo)
	if err != nil {
		return nil, err
	}

	return orderInfo, err
}
