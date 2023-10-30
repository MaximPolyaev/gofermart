package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

const pingTimeout = time.Second * 30

func InitDB(databaseURI string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %s", err)
	}

	return db, nil
}
