package entities

//OrderResponse - ответ системы расчета начислений баллов лояльности,
//Поля объекта ответа:
//	order — номер заказа;
//	status — статус расчёта начисления:
//		REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
//		INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
//		PROCESSING — расчёт начисления в процессе;
//		PROCESSED — расчёт начисления окончен;
//	accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
