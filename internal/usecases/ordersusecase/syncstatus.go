package ordersusecase

import (
	"context"
	"fmt"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/accrualstatus"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
)

const updateOrderAccrualsDelay = time.Millisecond * 500

func (uc *OrdersUseCase) StartSyncOrdersStatusesProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber := <-uc.updateOrdersCh:
			uc.updateOrderAccrualsWithDelay(ctx, orderNumber)
		}
	}
}

func (uc *OrdersUseCase) UpUpdateOrdersPool(ctx context.Context) error {
	orderNumbers, err := uc.storage.FindOrderNumbersToUpdateAccruals(ctx)
	if err != nil {
		return err
	}

	for _, orderNumber := range orderNumbers {
		uc.updateOrdersCh <- orderNumber
	}

	return nil
}

func (uc *OrdersUseCase) orderProcessWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber := <-uc.updateOrdersCh:
			uc.updateOrderAccrualsWithDelay(ctx, orderNumber)
		}
	}
}

func (uc *OrdersUseCase) updateOrderAccrualsWithDelay(ctx context.Context, orderNumber string) {
	tick := time.NewTicker(updateOrderAccrualsDelay)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			err := uc.updateOrderAccruals(ctx, orderNumber)
			if err != nil {
				uc.log.Error(err)

				uc.updateOrdersCh <- orderNumber
			}
			return
		}
	}
}

func (uc *OrdersUseCase) updateOrderAccruals(ctx context.Context, number string) error {
	err := uc.storage.ChangeOrderStatus(ctx, number, orderstatus.PROCESSING)
	if err != nil {
		return err
	}

	accrualOrder, err := uc.accrual.FetchAccrualOrder(ctx, number)
	if err != nil {
		return fmt.Errorf("fetch accrual %s: %s", number, err)
	}

	order := uc.orderStatusByAccrualStatus(accrualOrder)
	if order.Status == orderstatus.INVALID || order.Status == orderstatus.PROCESSED {
		err = uc.storage.SaveOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	if accrualOrder.IsNeedGetAccruals() {
		uc.updateOrdersCh <- number
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
