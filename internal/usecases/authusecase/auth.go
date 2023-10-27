package authusecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/utils/jwt"
)

type AuthUseCase struct {
	storage storage
}

type storage interface {
	FindUserWithPassword(ctx context.Context, login string) (*entities.UserWithPassword, error)
	CreateUser(ctx context.Context, user *entities.UserWithPassword) error
}

func New(storage storage) *AuthUseCase {
	return &AuthUseCase{
		storage: storage,
	}
}

func (a *AuthUseCase) ValidatePayload(payload entities.AuthPayload) error {
	if payload.Login == "" {
		return errors.New("login not be empty")
	}

	if payload.Password == "" {
		return errors.New("password not be empty")
	}

	return nil
}

func (a *AuthUseCase) SignIn(ctx context.Context, payload entities.AuthPayload) (string, error) {
	user, err := a.storage.FindUserWithPassword(ctx, payload.Login)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("incorrect login or password")
	}

	if a.hashPassword(payload.Password) != user.HashPassword {
		return "", errors.New("incorrect login or password")
	}

	return jwt.BuildToken(payload.Login)
}

func (a *AuthUseCase) SignUp(ctx context.Context, payload entities.AuthPayload) (string, error) {
	existUser, err := a.storage.FindUserWithPassword(ctx, payload.Login)
	if err != nil {
		return "", err
	}

	if existUser != nil {
		return "", fmt.Errorf("user %s is exist", payload.Login)
	}

	hashPassword := a.hashPassword(payload.Password)

	if err := a.storage.CreateUser(ctx, &entities.UserWithPassword{
		Login:        payload.Login,
		HashPassword: hashPassword,
	}); err != nil {
		return "", err
	}

	return jwt.BuildToken(payload.Login)
}

func (a *AuthUseCase) hashPassword(passwordStr string) string {
	h := sha256.New()
	h.Write([]byte(passwordStr))

	return hex.EncodeToString(h.Sum(nil))
}
