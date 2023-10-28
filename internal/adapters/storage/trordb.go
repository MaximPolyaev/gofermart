package storage

import (
	"context"
	"database/sql"

	"github.com/MaximPolyaev/gofermart/internal/utils/trmanager"
)

type conn interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)

	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
}

func (s *Storage) trOrDB(ctx context.Context) conn {
	txByCtx := ctx.Value(trmanager.DefaultTxKey)

	if txByCtx == nil {
		return s.db
	}

	tx, ok := txByCtx.(*sql.Tx)
	if ok {
		return tx
	}

	return s.db
}
