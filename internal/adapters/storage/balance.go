package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

func (s *Storage) LockUserWithCreateTx(ctx context.Context, userID int) (*sql.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `SELECT id FROM ref_user WHERE id = $1 FOR UPDATE`

	_, err = tx.ExecContext(ctx, q, userID)
	if err != nil {
		return nil, s.Rollback(tx, err)
	}

	return tx, nil
}

func (s *Storage) FindBalanceByUserID(ctx context.Context, tx *sql.Tx, userID int) (*entities.UserBalance, error) {
	q := `
SELECT sum(t.points) as current, sum(case when t.points < 0 then -1 * t.points else 0 end) as withdrawn
FROM reg_points_balance t
WHERE t.user_id = $1
`
	var balance entities.UserBalance
	var current sql.NullFloat64
	var withdrawn sql.NullFloat64
	var row *sql.Row

	if tx == nil {
		row = s.db.QueryRowContext(ctx, q, userID)
	} else {
		row = tx.QueryRowContext(ctx, q, userID)
	}

	err := row.Scan(&current, &withdrawn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &balance, nil
		}

		return &balance, err
	}

	if current.Valid {
		balance.Current = current.Float64
	}

	if withdrawn.Valid {
		balance.Withdrawn = withdrawn.Float64
	}

	return &balance, nil
}

func (s *Storage) FindWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error) {
	q := `
SELECT o.number, -1 * t.points, t.created_at
FROM reg_points_balance t
JOIN doc_order o ON o.id = t.order_id
WHERE o.user_id = $1 and t.points < 0
ORDER BY t.created_at
`

	rows, err := s.db.QueryContext(ctx, q, userID)
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				s.log.Error(err)
			}
		}
	}()
	if err != nil {
		return nil, err
	}

	var wroteOffs []entities.WroteOff

	for rows.Next() {
		var number string
		var sum float64
		var processedAt time.Time

		err := rows.Scan(&number, &sum, &processedAt)
		if err != nil {
			return nil, err
		}

		wroteOffs = append(wroteOffs, entities.WroteOff{
			WriteOff: entities.WriteOff{
				Order: number,
				Sum:   sum,
			},
			ProcessedAt: entities.RFC3339Time(processedAt),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wroteOffs, nil
}

func (s *Storage) WriteOffWithCommit(ctx context.Context, tx *sql.Tx, orderID int, userID int, points float64) error {
	q := `INSERT INTO reg_points_balance (order_id, user_id, points) VALUES ($1, $2, $3)`

	_, err := tx.ExecContext(ctx, q, orderID, userID, -1*points)
	if err != nil {
		return s.Rollback(tx, err)
	}

	return tx.Commit()
}
