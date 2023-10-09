package authusecase

import (
	"context"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type AuthUseCase struct{}

func New() *AuthUseCase {
	return &AuthUseCase{}
}

func (a *AuthUseCase) SignIn(ctx context.Context, payload entities.AuthPayload) (string, error) {
	return "", nil
}

func (a *AuthUseCase) SignUp(ctx context.Context, payload entities.AuthPayload) (string, error) {
	return "", nil
}
