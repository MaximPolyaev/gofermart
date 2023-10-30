package accrualusecase

import (
	"context"
	"errors"
	"github.com/MaximPolyaev/gofermart/internal/errors/accrualerrors"
	"testing"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/accrualstatus"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
	"github.com/MaximPolyaev/gofermart/internal/usecases/accrualusecase/mocks"
	"github.com/golang/mock/gomock"
)

func TestAccrualsUseCase_StartSyncOrdersStatusesProcess(t *testing.T) {
	type fetchedOrder struct {
		order *entities.AccrualOrder
		err   error
	}

	tests := []struct {
		name         string
		duration     time.Duration
		orderNumber  string
		fetchedOrder fetchedOrder
		updateOrder  *entities.Order
	}{
		{
			name:        "update accruals",
			duration:    time.Second * 2,
			orderNumber: "1",
			fetchedOrder: fetchedOrder{
				order: &entities.AccrualOrder{
					Order:   "1",
					Status:  accrualstatus.PROCESSED,
					Accrual: 100,
				},
			},
			updateOrder: &entities.Order{
				Number:  "1",
				Status:  orderstatus.PROCESSED,
				Accrual: 100,
			},
		},
		{
			name:        "update invalid",
			duration:    time.Second * 2,
			orderNumber: "1",
			fetchedOrder: fetchedOrder{
				order: &entities.AccrualOrder{
					Order:  "1",
					Status: accrualstatus.INVALID,
				},
			},
			updateOrder: &entities.Order{
				Number: "1",
				Status: orderstatus.INVALID,
			},
		},
		{
			name:        "fetch registered",
			duration:    time.Second * 2,
			orderNumber: "1",
			fetchedOrder: fetchedOrder{
				order: &entities.AccrualOrder{
					Order:  "1",
					Status: accrualstatus.REGISTERED,
				},
			},
		},
		{
			name:        "fetch processing",
			duration:    time.Second * 2,
			orderNumber: "1",
			fetchedOrder: fetchedOrder{
				order: &entities.AccrualOrder{
					Order:  "1",
					Status: accrualstatus.PROCESSING,
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.duration)
			defer cancel()

			accrual := mocks.NewMockaccrual(ctrl)
			accrual.EXPECT().FetchAccrualOrder(ctx, tt.orderNumber).Return(
				tt.fetchedOrder.order,
				tt.fetchedOrder.err,
			).AnyTimes()

			storage := mocks.NewMockstorage(ctrl)
			storage.EXPECT().FindOrderNumbersToUpdateAccruals(ctx).Return([]string{tt.orderNumber}, nil).AnyTimes()
			storage.EXPECT().ChangeOrderStatus(
				ctx,
				tt.orderNumber,
				orderstatus.PROCESSING,
			).AnyTimes()

			if tt.updateOrder != nil {
				storage.EXPECT().UpdateOrder(ctx, tt.updateOrder).Return(nil).AnyTimes()
			}

			log := mocks.NewMocklogger(ctrl)

			uc := New(accrual, storage, log)
			uc.StartSyncOrdersStatusesProcess(ctx)
		})
	}
}

func TestAccrualsUseCase_StartSyncOrdersStatusesProcess_RateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	orderNumber := "1"

	accrual := mocks.NewMockaccrual(ctrl)
	accrual.EXPECT().FetchAccrualOrder(ctx, orderNumber).Return(
		nil,
		&accrualerrors.RateLimitError{
			RetryAfter: 5,
		},
	)

	storage := mocks.NewMockstorage(ctrl)
	storage.EXPECT().FindOrderNumbersToUpdateAccruals(ctx).Return([]string{orderNumber}, nil).AnyTimes()
	storage.EXPECT().ChangeOrderStatus(
		ctx,
		orderNumber,
		orderstatus.PROCESSING,
	).AnyTimes()

	log := mocks.NewMocklogger(ctrl)
	log.EXPECT().Error(errors.New("update accruals 1: count reqs is many: over rate limit app"))
	log.EXPECT().Info("reset ticker to 5")

	uc := New(accrual, storage, log)
	uc.StartSyncOrdersStatusesProcess(ctx)
}
