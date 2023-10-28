package balanceusecase

import (
	"context"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/errors/balanceerrors"
)

type BalanceUseCase struct {
	storage storage
	trm     transactionManager
	log     logger
}

type storage interface {
	FindBalanceByUserID(ctx context.Context, userID int) (*entities.UserBalance, error)
	FindOrderIDByNumber(ctx context.Context, number string) (int, error)
	LockUserForUpdateBalance(ctx context.Context, userID int) error
	WriteOff(ctx context.Context, orderID int, userID int, points float64) error
	FindWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error)
}

type transactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type logger interface {
	Error(args ...interface{})
}

func New(storage storage, trm transactionManager, log logger) *BalanceUseCase {
	return &BalanceUseCase{
		storage: storage,
		trm:     trm,
		log:     log,
	}
}

func (uc *BalanceUseCase) GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error) {
	return uc.storage.FindBalanceByUserID(ctx, userID)
}

func (uc *BalanceUseCase) WriteOff(ctx context.Context, writeOff entities.WriteOff, userID int) error {
	orderID, err := uc.storage.FindOrderIDByNumber(ctx, writeOff.Order)
	if err != nil {
		return err
	}

	return uc.trm.Do(ctx, func(ctx context.Context) error {
		err := uc.storage.LockUserForUpdateBalance(ctx, userID)
		if err != nil {
			return err
		}

		balance, err := uc.GetBalance(ctx, userID)
		if err != nil {
			return err
		}

		if balance.Current <= 0 || balance.Current < writeOff.Sum {
			return balanceerrors.ErrInsufficientFunds
		}

		return uc.storage.WriteOff(ctx, orderID, userID, writeOff.Sum)
	})
}

func (uc *BalanceUseCase) GetWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error) {
	return uc.storage.FindWroteOffs(ctx, userID)
}
