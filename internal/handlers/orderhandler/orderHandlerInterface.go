package orderhandler

import "net/http"

var _ OrderHandler = (*orderHandlerImpl)(nil)

type OrderHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetOrdersHistory(writer http.ResponseWriter, request *http.Request)
}
