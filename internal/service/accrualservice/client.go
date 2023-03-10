package accrualservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sandor-clegane/go-market/internal/entities"
	"golang.org/x/time/rate"
)

const (
	getEndpoint = "api/orders"
	rpsLimit    = 1000
	duration    = 1 * time.Second
)

type AccrualService struct {
	accrualSystemAddress string
	httpClient           *http.Client
	rateLimiter          *rate.Limiter
}

func New(host string) *AccrualService {
	return &AccrualService{
		accrualSystemAddress: host,
		httpClient:           &http.Client{},
		rateLimiter:          rate.NewLimiter(rate.Every(duration), rpsLimit),
	}
}

func (as *AccrualService) do(method, endpoint string, number int) (*http.Response, error) {
	targetURL := fmt.Sprintf("%s/%s/%d", as.accrualSystemAddress, endpoint, number)

	req, err := http.NewRequest(method, targetURL, nil)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = as.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := as.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//GetOrderInfo - получение информации о расчёте начислений баллов лояльности.
func (as *AccrualService) GetOrderInfo(number int) (entities.OrderResponse, error) {
	resp, err := as.do(http.MethodGet, getEndpoint, number)
	if err != nil {
		return entities.OrderResponse{}, err
	}

	defer resp.Body.Close()
	var orderResponse entities.OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		return entities.OrderResponse{}, err
	}
	return orderResponse, nil
}
