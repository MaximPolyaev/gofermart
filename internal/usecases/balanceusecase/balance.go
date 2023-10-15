package balanceusecase

import (
	"context"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type BalanceUseCase struct {
	storage storage
}

type storage interface {
	FindBalanceByUserID(ctx context.Context, userID int) (*entities.UserBalance, error)
	FindBalanceByOrderNumber(ctx context.Context, number string) (float64, error)
	CreatePointsOperation(ctx context.Context, orderID int, points float64) error
	FindOrderIDByNumber(ctx context.Context, number string) (int, error)
	FindWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error)
}

func New(storage storage) *BalanceUseCase {
	return &BalanceUseCase{storage: storage}
}

func (uc *BalanceUseCase) GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error) {
	return uc.storage.FindBalanceByUserID(ctx, userID)
}

func (uc *BalanceUseCase) IsAvailableWriteOff(ctx context.Context, writeOff *entities.WriteOff) (bool, error) {
	balance, err := uc.storage.FindBalanceByOrderNumber(ctx, writeOff.Order)
	if err != nil {
		return false, err
	}

	return balance > 0 && balance >= writeOff.Sum, nil
}

func (uc *BalanceUseCase) WriteOff(ctx context.Context, off entities.WriteOff) error {
	orderID, err := uc.storage.FindOrderIDByNumber(ctx, off.Order)
	if err != nil {
		return err
	}

	return uc.storage.CreatePointsOperation(ctx, orderID, -1*off.Sum)
}

func (uc *BalanceUseCase) GetWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error) {
	return uc.storage.FindWroteOffs(ctx, userID)
}
