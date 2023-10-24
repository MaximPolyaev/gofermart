package ordersusecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type OrdersUseCase struct {
	storage storage
	accrual accrual
	log     logger
}

func New(storage storage, accrual accrual, log logger) *OrdersUseCase {
	return &OrdersUseCase{
		storage: storage,
		accrual: accrual,
		log:     log,
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
