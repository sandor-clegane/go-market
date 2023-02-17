package withdrawstorage

import (
	"context"
	"database/sql"

	"github.com/sandor-clegane/go-market/internal/entities"
)

const (
	insertWithdrawQuery = "" +
		"INSERT INTO withdraws (id, sum, processed_at, user_id) " +
		"VALUES ($1, $2, $3, $4)"
	getTotalWithdrawByUserIDQuery = "" +
		"SELECT SUM(sum) " +
		"FROM withdraws " +
		"WHERE user_id = $1"
	getAllWithdrawsByUserIDQuery = "" +
		"SELECT id, sum, " +
		"processed_at::timestamptz " +
		"FROM withdraws " +
		"WHERE user_id = $1"
)

type withdrawStorageImpl struct {
	db *sql.DB
}

func New(db *sql.DB) WithdrawStorage {
	return &withdrawStorageImpl{
		db: db,
	}

}

func (ws *withdrawStorageImpl) InsertWithdraw(ctx context.Context, withdraw entities.Withdraw) error {
	_, err := ws.db.ExecContext(ctx, insertWithdrawQuery,
		withdraw.Order, withdraw.Sum,
		withdraw.ProcessedAt, withdraw.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (ws *withdrawStorageImpl) GetTotalWithdrawnByUserID(ctx context.Context, userID string) (float32, error) {
	//It should be noted that except for count, these
	//aggregate functions in postgresql return a
	//null value when no rows are selected.
	row := ws.db.QueryRowContext(ctx, getTotalWithdrawByUserIDQuery, userID)
	var maybeTotalWithdraw sql.NullFloat64
	if err := row.Scan(&maybeTotalWithdraw); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
		return 0.0, nil
	}
	return float32(maybeTotalWithdraw.Float64), nil
}

func (ws *withdrawStorageImpl) GetAllWithdrawsByUserID(ctx context.Context, userID string) ([]entities.Withdraw, error) {
	rows, err := ws.db.QueryContext(ctx, getAllWithdrawsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Withdraw, 0)
	for rows.Next() {
		var w entities.Withdraw
		err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, w)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}
