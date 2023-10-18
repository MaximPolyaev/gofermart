package storage

import (
	"context"
	"database/sql"
)

type execCtx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
