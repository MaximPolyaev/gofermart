package storage

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) FindUserIDByOrderNumber(ctx context.Context, number int) (int, error) {
	var userId int

	q := `SELECT user_id FROM doc_order WHERE number = $1 LIMIT 1`

	err := s.db.QueryRowContext(ctx, q, number).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return userId, nil
}

func (s *Storage) CreateOrder(ctx context.Context, number int, userId int) error {
	q := `INSERT INTO doc_order (number, user_id) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, q, number, userId)

	return err
}
