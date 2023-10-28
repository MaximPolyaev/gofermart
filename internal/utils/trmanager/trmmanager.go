package trmanager

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	db *sql.DB
}

func New(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(context.WithValue(ctx, "tx", tx))
	if err != nil {
		txerr := tx.Rollback()
		if txerr != nil {
			return txerr
		}
		return err
	}

	return tx.Commit()
}
