package orderstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customerrors"
	"github.com/sandor-clegane/go-market/internal/entities/mappers"
	"github.com/sandor-clegane/go-market/internal/service/accrualservice"

	"github.com/omeid/pgerror"
)

const (
	workersCount   = 10
	recheckTimeout = 5 * time.Second
)

const (
	findOrderByNumberQuery = "" +
		"SELECT id, user_id FROM orders " +
		"WHERE id = $1"
	getAllOrdersByUserIDQuery = "" +
		"SELECT id, status, accrual_amount, " +
		"uploaded_at::timestamptz " +
		"FROM orders " +
		"WHERE user_id = $1"
	getTotalAccrualAmountByUserIDQuery = "" +
		"SELECT SUM(accrual_amount) " +
		"FROM orders " +
		"WHERE user_id = $1"
	getOrderNumbersWithNotFinalStatusQuery = "" +
		"SELECT id " +
		"FROM orders " +
		"WHERE status IN (1, 2)"
)

const (
	updateOrderByNumberQuery = "" +
		"UPDATE orders " +
		"SET status = $1, " +
		"accrual_amount = $2 " +
		"WHERE id = $3"
	insertOrderQuery = "" +
		"INSERT INTO orders (id, status, accrual_amount, uploaded_at, user_id) " +
		"VALUES ($1, $2, $3, $4, $5)"
)

type syncObject struct {
	wg   sync.WaitGroup
	once sync.Once
	done chan struct{}
}

type orderStorageImpl struct {
	db                 *sql.DB
	ticker             *time.Ticker
	updatingOrderQueue chan int
	so                 syncObject
	accrualService     *accrualservice.AccrualService
}

func New(db *sql.DB, accrualServiceAddress string) (OrderStorage, error) {
	storage := orderStorageImpl{
		so: syncObject{
			done: make(chan struct{}),
		},
		db:                 db,
		ticker:             time.NewTicker(recheckTimeout),
		updatingOrderQueue: make(chan int),
		accrualService:     accrualservice.New(accrualServiceAddress),
	}
	err := storage.runUpdateWorkersTaskScheduler()
	if err != nil {
		return nil, err
	}
	storage.runUpdatingStatusWorkerPool()
	return &storage, nil
}

func (os *orderStorageImpl) InsertOrder(ctx context.Context, order entities.Order) error {
	existedOrder, err := os.FindByNumber(ctx, order.Number)
	if err != nil {
		return err
	}
	if existedOrder.Number > 0 && existedOrder.UserID != order.UserID {
		return customerrors.NewExistedOrderError(existedOrder.Number, existedOrder.UserID)
	}
	_, err = os.db.ExecContext(ctx, insertOrderQuery,
		order.Number, order.Status, 0.0,
		order.UploadedAt, order.UserID)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return customerrors.NewOrderViolationError(order, err)
		}
	}
	return nil
}

func (os *orderStorageImpl) FindByNumber(ctx context.Context, number int) (entities.Order, error) {
	row := os.db.QueryRowContext(ctx, findOrderByNumberQuery, number)
	var existedOrder entities.Order
	err := row.Scan(&existedOrder.Number, &existedOrder.UserID)
	if err != nil && err != sql.ErrNoRows {
		return entities.Order{}, err
	}
	return existedOrder, nil
}

func (os *orderStorageImpl) GetAllOrdersByUserID(ctx context.Context, userID string) ([]entities.Order, error) {
	rows, err := os.db.QueryContext(ctx, getAllOrdersByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	orderList := make([]entities.Order, 0)
	for rows.Next() {
		var o entities.Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return orderList, nil
}

func (os *orderStorageImpl) GetTotalAccrualAmountByUserID(ctx context.Context, userID string) (float32, error) {
	var totalAccrualAmount float32
	row := os.db.QueryRowContext(ctx, getTotalAccrualAmountByUserIDQuery, userID)
	err := row.Scan(&totalAccrualAmount)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return totalAccrualAmount, nil
}

func (os *orderStorageImpl) getOrderNumbersWithNotFinalStatus(ctx context.Context) ([]int, error) {
	rows, err := os.db.QueryContext(ctx, getOrderNumbersWithNotFinalStatusQuery)
	if err != nil {
		return nil, err
	}

	nums := make([]int, 0)
	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return nums, nil
}

func (os *orderStorageImpl) addNumberToUpdatingQueue(number int) error {
	select {
	case <-os.so.done:
		return fmt.Errorf("can`t add %d number as storage stopped", number)
	case os.updatingOrderQueue <- number:
		return nil
	}
}

func (os *orderStorageImpl) runUpdateWorkersTaskScheduler() error {
	var err error
	go func() {
		for {
			select {
			case <-os.so.done:
				return
			case <-os.ticker.C:
				nums, err := os.getOrderNumbersWithNotFinalStatus(context.Background())
				if err != nil {
					return
				}
				for _, num := range nums {
					err = os.addNumberToUpdatingQueue(num)
					if err != nil {
						return
					}
				}
			}
		}
	}()
	return err
}

func (os *orderStorageImpl) runUpdatingStatusWorkerPool() {
	os.so.wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go func(workerID int) {
			for num := range os.updatingOrderQueue {
				orderResponse, err := os.accrualService.GetOrderInfo(num)
				if err != nil {
					log.Printf("get order response error %v", err)
					continue
				}
				order, err := mappers.MapOrderResponseToOrder(orderResponse)
				if err != nil {
					log.Printf("mapper error %v", err)
					continue
				}
				_, err = os.db.ExecContext(context.Background(), updateOrderByNumberQuery,
					order.Status, order.Accrual, order.Number)
				if err != nil {
					log.Printf("update order error %v", err)
					continue
				}
			}
			os.so.wg.Done()
			log.Printf("Worker %d done", workerID)
		}(i)
	}
}

func (os *orderStorageImpl) StopSchedulerAndWorkerPool() {
	os.so.once.Do(func() {
		close(os.so.done)
		close(os.updatingOrderQueue)
	})
	os.so.wg.Wait()
	os.ticker.Stop()
}
