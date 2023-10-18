package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
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

func (s *Storage) SaveOrder(ctx context.Context, order *entities.Order) error {
	if order.Accrual == 0 {
		return s.ChangeOrderStatus(ctx, order.Number, order.Status)
	}

	orderID, err := s.FindOrderIDByNumber(ctx, order.Number)
	if err != nil {
		return err
	}

	userID, err := s.FindUserIDByOrderNumber(ctx, order.Number)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = s.changeOrderStatus(ctx, tx, order.Number, order.Status)
	if err != nil {
		return s.rollback(tx, err)
	}

	err = s.createPointsOperation(ctx, tx, orderID, userID, order.Accrual)
	if err != nil {
		return s.rollback(tx, err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
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
		var status orderstatus.OrderStatus
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

func (s *Storage) FindOrderNumbersToUpdateAccruals(ctx context.Context) ([]string, error) {
	q := `
SELECT number
FROM doc_order
WHERE status IN ($1, $2)
`

	rows, err := s.db.QueryContext(ctx, q, orderstatus.NEW, orderstatus.PROCESSING)
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

	var orders []string

	for rows.Next() {
		var number string

		err := rows.Scan(&number)
		if err != nil {
			return nil, err
		}

		orders = append(orders, number)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) ChangeOrderStatus(ctx context.Context, number string, status orderstatus.OrderStatus) error {
	return s.changeOrderStatus(ctx, s.db, number, status)
}

func (s *Storage) changeOrderStatus(
	ctx context.Context,
	ex execCtx,
	number string,
	status orderstatus.OrderStatus,
) error {
	q := `UPDATE doc_order SET status = $1, changed_at = now() WHERE number = $2`

	_, err := ex.ExecContext(ctx, q, status, number)

	return err
}

func (s *Storage) rollback(tx *sql.Tx, err error) error {
	txErr := tx.Rollback()
	if txErr != nil {
		return txErr
	}

	return err
}
