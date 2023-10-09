package storage

import (
	"context"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

func (s *Storage) GetUserWithPassword(ctx context.Context, login string) (*entities.UserWithPassword, error) {
	return nil, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *entities.UserWithPassword) error {
	return nil
}
