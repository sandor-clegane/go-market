package withdrawhandler

import (
	"net/http"
)

var _ WithdrawHandler = (*withdrawHandlerImpl)(nil)

type WithdrawHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetUserBalance(writer http.ResponseWriter, request *http.Request)
	GetWithdrawalsHistory(writer http.ResponseWriter, request *http.Request)
}
