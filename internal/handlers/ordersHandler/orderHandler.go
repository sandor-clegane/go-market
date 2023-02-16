package ordersHandler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sandor-clegane/go-market/internal/entities/customErrors"
	"github.com/sandor-clegane/go-market/internal/service/cookieService"
	"github.com/sandor-clegane/go-market/internal/service/orderService"
)

type orderHandlerImpl struct {
	orderService  orderService.OrderService
	cookieService cookieService.CookieService
}

func New(orderService orderService.OrderService, cookieService cookieService.CookieService) OrderHandler {
	return &orderHandlerImpl{orderService, cookieService}
}

//Create Хендлер: POST /api/user/orders.
//Хендлер доступен только аутентифицированным пользователям.
//Номером заказа является последовательность цифр произвольной длины.
//Формат запроса:
//
//POST /api/user/orders HTTP/1.1
//Content-Type: text/plain
//...
//
//12345678903
func (o *orderHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	userID := o.cookieService.AuthenticateUser(writer, request)

	orderNumber, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err = o.orderService.CreateOrder(request.Context(), string(orderNumber), userID)

	if err != nil {
		// номер заказа уже был загружен этим пользователем;
		var ov *customErrors.OrderViolationError
		if errors.As(err, &ov) {
			http.Error(writer, err.Error(), http.StatusOK)
			return
		}
		// номер заказа уже был загружен другим пользователем;
		var eo *customErrors.ExistedOrderError
		if errors.As(err, &eo) {
			http.Error(writer, err.Error(), http.StatusConflict)
			return
		}
		// неверный формат номера заказа;
		var iof *customErrors.InvalidOrderNumberFormatError
		if errors.As(err, &iof) {
			http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusAccepted)
}

//GetAll Хендлер: GET /api/user/orders.
//Хендлер доступен только авторизованному пользователю.
//Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым.
//Формат даты — RFC3339.
//Доступные статусы обработки расчётов:
//	NEW — заказ загружен в систему, но не попал в обработку;
//	PROCESSING — вознаграждение за заказ рассчитывается;
//	INVALID — система расчёта вознаграждений отказала в расчёте;
//	PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
//Формат запроса:
//
//GET /api/user/orders HTTP/1.1
//Content-Length: 0
func (o *orderHandlerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID := o.cookieService.AuthenticateUser(writer, request)

	ordersListSorted, notFoundErr := o.orderService.GetAllOrdersByUserID(request.Context(), userID)
	if notFoundErr != nil {
		http.Error(writer, notFoundErr.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writeErr := json.NewEncoder(writer).Encode(ordersListSorted)

	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
