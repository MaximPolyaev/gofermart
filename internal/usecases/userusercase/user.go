package userusercase

import "context"

type UserUseCase struct {
	storage storage
}

type storage interface {
	FindUserIDByLogin(ctx context.Context, login string) (int, error)
}

func New(storage storage) *UserUseCase {
	return &UserUseCase{storage: storage}
}

func (uc *UserUseCase) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	return uc.storage.FindUserIDByLogin(ctx, login)
}
