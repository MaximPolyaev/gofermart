package balanceusecase

import (
	"context"
	"testing"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/errors/balanceerrors"
	"github.com/MaximPolyaev/gofermart/internal/usecases/balanceusecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockTrm struct{}

func (trm *MockTrm) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestBalanceUseCase_WriteOff(t *testing.T) {
	tests := []struct {
		name                string
		wantInsufficientErr bool
	}{
		{
			name:                "invalid balance",
			wantInsufficientErr: true,
		},
		{
			name:                "valid balance",
			wantInsufficientErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMocklogger(ctrl)
	trm := &MockTrm{}
	ctx := context.TODO()
	userID := 1
	orderID := 1
	orderNumber := "1"
	writeOff := entities.WriteOff{
		Order: orderNumber,
		Sum:   100,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockstorage(ctrl)
			store.EXPECT().FindOrderIDByNumber(ctx, orderNumber).Return(orderID, nil)
			store.EXPECT().LockUserForUpdateBalance(ctx, userID)

			var balance float64

			if !tt.wantInsufficientErr {
				balance = 200
				store.EXPECT().WriteOff(ctx, orderID, userID, writeOff.Sum).Return(nil)
			}

			store.EXPECT().FindBalanceByUserID(ctx, userID).Return(&entities.UserBalance{
				Current: balance,
			}, nil)

			uc := New(store, trm, log)

			err := uc.WriteOff(ctx, writeOff, userID)

			if tt.wantInsufficientErr {
				assert.EqualError(t, err, balanceerrors.ErrInsufficientFunds.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
