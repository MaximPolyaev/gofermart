package balanceusecase

import (
	"context"
	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type BalanceUseCase struct {
	storage storage
}

type storage interface {
	FindBalanceByUserId(ctx context.Context, userID int) (*entities.UserBalance, error)
}

func New(storage storage) *BalanceUseCase {
	return &BalanceUseCase{storage: storage}
}

func (uc *BalanceUseCase) GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error) {
	return uc.storage.FindBalanceByUserId(ctx, userID)
}
