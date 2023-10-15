package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums"
)

func (s *Storage) FindUserIDByOrderNumber(ctx context.Context, number string) (int, error) {
	var userID int

	q := `SELECT user_id FROM doc_order WHERE number = $1 LIMIT 1`

	err := s.db.QueryRowContext(ctx, q, number).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return userID, nil
}

func (s *Storage) CreateOrder(ctx context.Context, number string, userID int) error {
	q := `INSERT INTO doc_order (number, user_id) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, q, number, userID)

	return err
}

func (s *Storage) FindOrdersByUserID(ctx context.Context, userID int) ([]entities.Order, error) {
	q := `
SELECT t.number, t.status, t.created_at, coalesce(b.points, 0) as points
FROM doc_order t
LEFT JOIN (
    SELECT t.order_id, sum(t.points) points
    FROM reg_points_balance t
    GROUP BY t.order_id
) b ON b.order_id = t.id
WHERE t.user_id = $1
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

	var orders []entities.Order

	for rows.Next() {
		var number string
		var status enums.OrderStatus
		var createdAt time.Time
		var points float64

		err := rows.Scan(&number, &status, &createdAt, &points)
		if err != nil {
			return nil, err
		}

		orders = append(orders, entities.Order{
			Number:     number,
			Status:     status,
			Accrual:    points,
			UploadedAt: entities.RFC3339Time(createdAt),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) FindOrderIDByNumber(ctx context.Context, number string) (int, error) {
	var id int

	q := `SELECT id FROM doc_order WHERE number = $1 LIMIT 1`

	err := s.db.QueryRowContext(ctx, q, number).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
