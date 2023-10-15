package ordersusecase

import (
	"context"
	"errors"
	"strconv"
)

type OrdersUseCase struct {
	storage storage
}

type storage interface {
	FindUserIDByOrderNumber(ctx context.Context, number string) (int, error)
	CreateOrder(ctx context.Context, number string, userID int) error
}

func New(storage storage) *OrdersUseCase {
	return &OrdersUseCase{storage: storage}
}

func (us *OrdersUseCase) GetUserID(ctx context.Context, number string) (int, error) {
	return us.storage.FindUserIDByOrderNumber(ctx, number)
}

func (us *OrdersUseCase) CreateOrder(ctx context.Context, number string, userID int) error {
	return us.storage.CreateOrder(ctx, number, userID)
}

func (us *OrdersUseCase) ValidateNumber(number string) error {
	sum, err := us.getLuhnSum(number)
	if err != nil {
		return err
	}

	if sum%10 != 0 {
		return errors.New("invalid order number")
	}

	return nil
}

func (us *OrdersUseCase) getLuhnSum(number string) (int64, error) {
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
