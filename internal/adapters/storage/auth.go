package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

func (s *Storage) FindUserWithPassword(ctx context.Context, login string) (*entities.UserWithPassword, error) {
	var user entities.UserWithPassword

	q := `SELECT login, password FROM ref_user WHERE login = $1 LIMIT 1`

	err := s.db.QueryRowContext(ctx, q, login).Scan(&user.Login, &user.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *entities.UserWithPassword) error {
	q := `
		INSERT INTO ref_user (login, password)
		VALUES ($1, $2)
	`

	if _, err := s.db.ExecContext(ctx, q, user.Login, user.HashPassword); err != nil {
		return err
	}

	return nil
}
