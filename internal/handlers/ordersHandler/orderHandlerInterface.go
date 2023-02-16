package ordersHandler

import "net/http"

var _ OrderHandler = &orderHandlerImpl{}

type OrderHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
}
