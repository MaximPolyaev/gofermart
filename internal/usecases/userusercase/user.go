package userusercase

import "context"

type UserUseCase struct {
	storage storage
}

type storage interface {
	FindUserIdByLogin(ctx context.Context, login string) (int, error)
}

func New(storage storage) *UserUseCase {
	return &UserUseCase{storage: storage}
}

func (uc *UserUseCase) GetUserIdByLogin(ctx context.Context, login string) (int, error) {
	return uc.storage.FindUserIdByLogin(ctx, login)
}
