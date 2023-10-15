package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

func (s *Storage) FindBalanceByUserId(ctx context.Context, userID int) (*entities.UserBalance, error) {
	var balance entities.UserBalance

	q := `
SELECT sum(t.points) as current, sum(case when t.points < 0 then -1 * t.points else 0 end) as withdrawn
FROM reg_points_balance t
JOIN doc_order o ON o.id = t.order_id
WHERE o.user_id = $1
`

	err := s.db.QueryRowContext(ctx, q, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &balance, nil
		}

		return &balance, err
	}

	return &balance, nil
}
