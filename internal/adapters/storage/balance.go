package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

func (s *Storage) FindBalanceByUserID(ctx context.Context, userID int) (*entities.UserBalance, error) {
	var balance entities.UserBalance

	q := `
SELECT sum(t.points) as current, sum(case when t.points < 0 then -1 * t.points else 0 end) as withdrawn
FROM reg_points_balance t
JOIN doc_order o ON o.id = t.order_id
WHERE o.user_id = $1
`

	var current sql.NullFloat64
	var withdrawn sql.NullFloat64

	err := s.db.QueryRowContext(ctx, q, userID).Scan(&current, &withdrawn)
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

func (s *Storage) FindBalanceByOrderNumber(ctx context.Context, number string) (float64, error) {
	q := `
SELECT sum(t.points) 
FROM reg_points_balance t 
JOIN doc_order o ON o.id = t.order_id
WHERE o.number = $1
    `

	var points sql.NullFloat64

	err := s.db.QueryRowContext(ctx, q, number).Scan(&points)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return points.Float64, nil
}

func (s *Storage) CreatePointsOperation(ctx context.Context, orderID int, points float64) error {
	return s.createPointsOperation(ctx, s.db, orderID, points)
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
				fmt.Println("rows close", err)
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

func (s *Storage) createPointsOperation(ctx context.Context, ex execCtx, orderID int, points float64) error {
	q := `INSERT INTO reg_points_balance (order_id, points) VALUES ($1, $2)`

	_, err := ex.ExecContext(ctx, q, orderID, points)
	if err != nil {
		return err
	}

	return nil
}
