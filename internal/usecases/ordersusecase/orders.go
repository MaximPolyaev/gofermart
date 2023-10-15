package ordersusecase

import "context"

type OrdersUseCase struct {
	storage storage
}

type storage interface {
	FindUserIDByOrderNumber(ctx context.Context, number int) (int, error)
	CreateOrder(ctx context.Context, number int, userId int) error
}

func New(storage storage) *OrdersUseCase {
	return &OrdersUseCase{storage: storage}
}

func (us *OrdersUseCase) GetUserID(ctx context.Context, number int) (int, error) {
	return us.storage.FindUserIDByOrderNumber(ctx, number)
}

func (us *OrdersUseCase) CreateOrder(ctx context.Context, number int, userId int) error {
	return us.storage.CreateOrder(ctx, number, userId)
}

func (us *OrdersUseCase) ValidateLuhn(number int) bool {
	return us.getLuhnSum(number)%10 == 0
}

func (us *OrdersUseCase) getLuhnSum(number int) int {
	var sum int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur -= 9
			}
		}

		sum += cur

		number = number / 10
	}

	return sum % 10
}
