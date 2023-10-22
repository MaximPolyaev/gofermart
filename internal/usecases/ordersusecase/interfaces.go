package ordersusecase

import (
	"context"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
)

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

type logger interface {
	Error(args ...interface{})
}
