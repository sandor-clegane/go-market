package withdrawHandler

import (
	"net/http"
)

var _ WithdrawHandler = &withdrawHandlerImpl{}

type WithdrawHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetUserBalance(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
}
