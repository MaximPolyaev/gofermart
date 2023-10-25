package ordersusecase

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

func (uc *OrdersUseCase) StartSyncOrdersStatusesProcess(ctx context.Context) {
	tickerDuration := time.Second

	tick := time.NewTicker(tickerDuration)
	defer tick.Stop()

	var orderNumbers []string

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
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
			if err != nil {
				if len(orderNumbers) == 1 {
					orderNumbers = make([]string, 0)
					continue
				}
				orderNumbers = orderNumbers[1:]
				continue
			}

			uc.log.Error(fmt.Errorf("update accruals %s: %s", orderNumber, err))

			if errors.Is(err, accrualerrors.ErrRateLimit) {
				tickerDuration *= 2
				tick.Reset(tickerDuration)
				uc.log.Info(fmt.Sprintf("increase ticker to %d", tickerDuration/time.Second))
			}
		}
	}
}

func (uc *OrdersUseCase) updateOrderAccruals(ctx context.Context, number string) error {
	err := uc.storage.ChangeOrderStatus(ctx, number, orderstatus.PROCESSING, nil)
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
