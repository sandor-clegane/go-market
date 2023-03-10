package mappers

import (
	"strconv"
	"time"

	"github.com/sandor-clegane/go-market/internal/entities"
)

func MapToWithdraw(withdrawRequest entities.WithdrawRequest, orderNumber int, userID string) entities.Withdraw {
	return entities.Withdraw{
		Order:       orderNumber,
		Sum:         withdrawRequest.Sum,
		ProcessedAt: time.Now(),
		UserID:      userID,
	}
}

func MapToBalanceRequest(totalUserAccrual float32, totalUserWithdraw float32) entities.BalanceRequest {
	return entities.BalanceRequest{
		Current:   totalUserAccrual - totalUserWithdraw,
		Withdrawn: totalUserWithdraw,
	}
}

func MapWithdrawListToWithdrawDTOList(withdrawList []entities.Withdraw) []entities.WithdrawDTO {
	withdrawDTOList := make([]entities.WithdrawDTO, 0, len(withdrawList))
	var dto entities.WithdrawDTO
	for _, w := range withdrawList {
		dto = entities.WithdrawDTO{
			Order:       strconv.Itoa(w.Order),
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
		}
		withdrawDTOList = append(withdrawDTOList, dto)
	}
	return withdrawDTOList
}
