package accrualusecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/accrualstatus"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
	"github.com/MaximPolyaev/gofermart/internal/errors/accrualerrors"
)

const fetchAccrualsInterval = time.Second

type AccrualsUseCase struct {
	accrual accrual
	storage storage
	log     logger
}

type accrual interface {
	FetchAccrualOrder(ctx context.Context, number string) (*entities.AccrualOrder, error)
}

type logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

type storage interface {
	FindOrderNumbersToUpdateAccruals(ctx context.Context) ([]string, error)
	ChangeOrderStatus(ctx context.Context, number string, status orderstatus.OrderStatus) error
	UpdateOrder(ctx context.Context, order *entities.Order) error
}

func New(accrual accrual, storage storage, log logger) *AccrualsUseCase {
	return &AccrualsUseCase{
		accrual: accrual,
		storage: storage,
		log:     log,
	}
}

func (uc *AccrualsUseCase) StartSyncOrdersStatusesProcess(ctx context.Context) {
	tick := time.NewTicker(fetchAccrualsInterval)
	defer tick.Stop()

	var orderNumbers []string
	var isErrTick bool

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if isErrTick {
				isErrTick = false
				tick.Reset(fetchAccrualsInterval)
			}

			var err error

			if len(orderNumbers) == 0 {
				orderNumbers, err = uc.storage.FindOrderNumbersToUpdateAccruals(ctx)
				if err != nil {
					uc.log.Error(err)
					continue
				}

				if len(orderNumbers) == 0 {
					continue
				}
			}

			orderNumber := orderNumbers[0]

			err = uc.updateOrderAccruals(ctx, orderNumber)
			if err == nil {
				if len(orderNumbers) == 1 {
					orderNumbers = make([]string, 0)
					continue
				}
				orderNumbers = orderNumbers[1:]
				continue
			}

			uc.log.Error(fmt.Errorf("update accruals %s: %s", orderNumber, err))

			var rateLimitErr *accrualerrors.RateLimitError
			if errors.As(err, &rateLimitErr) {
				tick.Reset(time.Duration(rateLimitErr.RetryAfter) * time.Second)
				isErrTick = true
				uc.log.Info(fmt.Sprintf("reset ticker to %d", rateLimitErr.RetryAfter))
			}
		}
	}
}

func (uc *AccrualsUseCase) updateOrderAccruals(ctx context.Context, number string) error {
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
		err = uc.storage.UpdateOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *AccrualsUseCase) orderStatusByAccrualStatus(accrualOrder *entities.AccrualOrder) *entities.Order {
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
