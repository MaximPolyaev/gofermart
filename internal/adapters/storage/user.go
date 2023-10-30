package storage

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) FindUserIDByLogin(ctx context.Context, login string) (int, error) {
	var id int

	q := `SELECT id FROM ref_user WHERE login = $1 LIMIT 1`

	err := s.trOrDB(ctx).QueryRowContext(ctx, q, login).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return id, nil
}
