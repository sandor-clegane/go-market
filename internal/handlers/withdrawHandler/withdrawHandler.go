package withdrawHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customErrors"
	"github.com/sandor-clegane/go-market/internal/service/cookieService"
	"github.com/sandor-clegane/go-market/internal/service/withdrawService"
)

type withdrawHandlerImpl struct {
	withdrawService withdrawService.WithdrawService
	cookieService   cookieService.CookieService
}

func New(withdrawService withdrawService.WithdrawService,
	cookieService cookieService.CookieService) WithdrawHandler {
	return &withdrawHandlerImpl{
		withdrawService,
		cookieService,
	}
}

//Create Хендлер: POST /api/user/balance/withdraw
//Хендлер доступен только авторизованному пользователю.
//Номер заказа представляет собой гипотетический номер нового заказа
//пользователя, в счёт оплаты которого списываются баллы.
//Формат запроса:
//
//POST /api/user/balance/withdraw HTTP/1.1
//Content-Type: application/json
//
//{
//	"order": "2377225624",
//	"sum": 751
//}
func (w *withdrawHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	userID := w.cookieService.AuthenticateUser(writer, request)

	var req entities.WithdrawRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err = w.withdrawService.CreateWithdraw(request.Context(), req, userID)

	if err != nil {
		//на счету недостаточно средств;
		var le *customErrors.LimitExceededError
		if errors.As(err, &le) {
			http.Error(writer, err.Error(), http.StatusPaymentRequired)
			return
		}
		// неверный номер заказа;
		var iof *customErrors.InvalidOrderNumberFormatError
		if errors.As(err, &iof) {
			http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

//GetUserBalance Хендлер: GET /api/user/balance.
//Хендлер доступен только авторизованному пользователю.
//В ответе должны содержаться данные о текущей сумме баллов лояльности,
//а также сумме использованных за весь период регистрации баллов.
//Формат запроса:
//
//GET /api/user/balance HTTP/1.1
//Content-Length: 0
func (w *withdrawHandlerImpl) GetUserBalance(writer http.ResponseWriter, request *http.Request) {
	userID := w.cookieService.AuthenticateUser(writer, request)

	balanceRequest, err := w.withdrawService.GetBalanceInfoByID(request.Context(), userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writeErr := json.NewEncoder(writer).Encode(balanceRequest)
	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}

//GetAll Хендлер: GET /api/user/withdrawals.
//Хендлер доступен только авторизованному пользователю.
//Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым.
//Формат даты — RFC3339.
//Формат запроса:
//
//GET /api/user/withdrawals HTTP/1.1
//Content-Length: 0
//
//Формат ответа:
//
//200 OK HTTP/1.1
//Content-Type: application/json
//...
//
//[
//	{
//		"order": "2377225624",
//		"sum": 500,
//		"processed_at": "2020-12-09T16:09:57+03:00"
//	}
//]
func (w *withdrawHandlerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID := w.cookieService.AuthenticateUser(writer, request)

	withdrawDTOListSorted, err := w.withdrawService.GetWithdrawsInfoByID(request.Context(), userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if len(withdrawDTOListSorted) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(withdrawDTOListSorted)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
