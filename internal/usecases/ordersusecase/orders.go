package ordersusecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/accrualstatus"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
	"strconv"
	"sync"
	"time"
)

const processOrdersWorkersCount = 3
const updateOrderAccrualsDelay = time.Millisecond * 500

type OrdersUseCase struct {
	storage        storage
	accrual        accrual
	updateOrdersCh chan string
}

type storage interface {
	FindUserIDByOrderNumber(ctx context.Context, number string) (int, error)
	CreateOrder(ctx context.Context, number string, userID int) error
	SaveOrder(ctx context.Context, order *entities.Order) error
	FindOrdersByUserID(ctx context.Context, userID int) ([]entities.Order, error)
	FindOrderNumbersToUpdateAccruals(ctx context.Context) ([]string, error)
	ChangeOrderStatus(ctx context.Context, number string, status orderstatus.OrderStatus) error
}

type accrual interface {
	FetchAccrualOrder(ctx context.Context, number string) (*entities.AccrualOrder, error)
}

func New(storage storage, accrual accrual) *OrdersUseCase {
	return &OrdersUseCase{
		storage:        storage,
		accrual:        accrual,
		updateOrdersCh: make(chan string, processOrdersWorkersCount),
	}
}

func (uc *OrdersUseCase) GetUserID(ctx context.Context, number string) (int, error) {
	return uc.storage.FindUserIDByOrderNumber(ctx, number)
}

func (uc *OrdersUseCase) CreateOrder(ctx context.Context, number string, userID int) error {
	err := uc.storage.CreateOrder(ctx, number, userID)
	if err != nil {
		return err
	}

	uc.addOrderToUpdateAccruals(number)

	return nil
}

func (uc *OrdersUseCase) GetOrders(ctx context.Context, userID int) ([]entities.Order, error) {
	return uc.storage.FindOrdersByUserID(ctx, userID)
}

func (uc *OrdersUseCase) ValidateNumber(number string) error {
	if number == "" {
		return errors.New("empty number")
	}

	sum, err := uc.getLuhnSum(number)
	if err != nil {
		return err
	}

	if sum%10 != 0 {
		return errors.New("invalid order number")
	}

	return nil
}

func (uc *OrdersUseCase) getLuhnSum(number string) (int64, error) {
	var sum int64

	dOnIncrease := len(number) % 2

	for i, dRune := range number {
		d, err := strconv.Atoi(string(dRune))
		if err != nil {
			return 0, err
		}

		if i%2 == dOnIncrease {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += int64(d)
	}

	return sum, nil
}

func (uc *OrdersUseCase) StartUpdateOrdersProcess(ctx context.Context) error {
	err := uc.upUpdateOrdersPool(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < processOrdersWorkersCount; i++ {
		uc.makeOrderProcessWorker(ctx)
	}

	return nil
}

func (uc *OrdersUseCase) upUpdateOrdersPool(ctx context.Context) error {
	orderNumbers, err := uc.storage.FindOrderNumbersToUpdateAccruals(ctx)
	if err != nil {
		return err
	}

	for _, orderNumber := range orderNumbers {
		uc.addOrderToUpdateAccruals(orderNumber)
	}

	return nil
}

func (uc *OrdersUseCase) addOrderToUpdateAccruals(number string) {
	go func() {
		uc.updateOrdersCh <- number
	}()
}

func (uc *OrdersUseCase) makeOrderProcessWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case orderNumber := <-uc.updateOrdersCh:
				uc.updateOrderAccrualsWithDelay(ctx, orderNumber)
			}
		}
	}()
}

func (uc *OrdersUseCase) updateOrderAccrualsWithDelay(ctx context.Context, orderNumber string) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		tick := time.NewTicker(updateOrderAccrualsDelay)
		defer func() {
			tick.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				err := uc.updateOrderAccruals(ctx, orderNumber)
				if err != nil {
					uc.addOrderToUpdateAccruals(orderNumber)
					fmt.Println("update order accruals", err)
				}
				return
			}
		}
	}()

	wg.Wait()
}

func (uc *OrdersUseCase) updateOrderAccruals(ctx context.Context, number string) error {
	err := uc.storage.ChangeOrderStatus(ctx, number, orderstatus.PROCESSING)
	if err != nil {
		return err
	}

	accrualOrder, err := uc.accrual.FetchAccrualOrder(ctx, number)
	if err != nil {
		return err
	}

	order := uc.orderStatusByAccrualStatus(accrualOrder)
	if order.Status == orderstatus.INVALID || order.Status == orderstatus.PROCESSED {
		err = uc.storage.SaveOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	if accrualOrder.IsNeedGetAccruals() {
		uc.addOrderToUpdateAccruals(number)
	}

	return nil
}

func (uc *OrdersUseCase) orderStatusByAccrualStatus(accrualOrder *entities.AccrualOrder) *entities.Order {
	var order entities.Order

	order.Number = accrualOrder.Order
	order.Accrual = accrualOrder.Accrual

	switch accrualOrder.Status {
	case accrualstatus.REGISTERED:
		order.Status = orderstatus.NEW
	case accrualstatus.INVALID:
		order.Status = orderstatus.INVALID
	case accrualstatus.PROCESSING:
		order.Status = orderstatus.PROCESSING
	case accrualstatus.PROCESSED:
		order.Status = orderstatus.PROCESSED
	}

	return &order
}
