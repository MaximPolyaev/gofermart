package balanceusecase

import (
	"context"
	"database/sql"
	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/errors/balanceerrors"
)

type BalanceUseCase struct {
	storage storage
	log     logger
}

type storage interface {
	FindBalanceByUserID(ctx context.Context, tx *sql.Tx, userID int) (*entities.UserBalance, error)
	FindOrderIDByNumber(ctx context.Context, number string) (int, error)
	WriteOffWithCommit(ctx context.Context, tx *sql.Tx, orderID int, userID int, points float64) error
	FindWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error)
	LockUserWithCreateTx(ctx context.Context, userID int) (*sql.Tx, error)
	Rollback(tx *sql.Tx, err error) error
}

type logger interface {
	Error(args ...interface{})
}

func New(storage storage, log logger) *BalanceUseCase {
	return &BalanceUseCase{storage: storage, log: log}
}

func (uc *BalanceUseCase) GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error) {
	return uc.storage.FindBalanceByUserID(ctx, nil, userID)
}

func (uc *BalanceUseCase) WriteOff(ctx context.Context, writeOff entities.WriteOff, userID int) error {
	orderID, err := uc.storage.FindOrderIDByNumber(ctx, writeOff.Order)
	if err != nil {
		return err
	}

	tx, err := uc.storage.LockUserWithCreateTx(ctx, userID)
	if err != nil {
		return err
	}

	balance, err := uc.storage.FindBalanceByUserID(ctx, tx, userID)
	if err != nil {
		return uc.storage.Rollback(tx, err)
	}

	if balance.Current <= 0 || balance.Current < writeOff.Sum {
		return uc.storage.Rollback(tx, balanceerrors.ErrInsufficientFunds)
	}

	return uc.storage.WriteOffWithCommit(ctx, tx, orderID, userID, writeOff.Sum)
}

func (uc *BalanceUseCase) GetWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error) {
	return uc.storage.FindWroteOffs(ctx, userID)
}
