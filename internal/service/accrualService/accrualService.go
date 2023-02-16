package accrualService

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sandor-clegane/go-market/internal/entities"
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
func (ac *AccrualService) GetOrderInfo(orderID int) (entities.OrderResponse, error) {
	targetURL := fmt.Sprintf("%s/%s/%d", ac.accrualSystemAddress, getEndpoint, orderID)

	res, err := http.Get(targetURL)
	if err != nil {
		return entities.OrderResponse{}, err
	}
	defer res.Body.Close()
	var orderInfo entities.OrderResponse
	err = json.NewDecoder(res.Body).Decode(&orderInfo)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	return orderInfo, err
}
