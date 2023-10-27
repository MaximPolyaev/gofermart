package authusecase

import (
	"context"
	"errors"
	"testing"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/usecases/authusecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const stubHashPassword = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"

func TestAuthUseCase_ValidatePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload entities.AuthPayload
		wantErr bool
	}{
		{
			name:    "login not be empty",
			payload: entities.AuthPayload{},
			wantErr: true,
		},
		{
			name:    "password not be empty",
			payload: entities.AuthPayload{Login: "login"},
			wantErr: true,
		},
		{
			name:    "without error",
			payload: entities.AuthPayload{Login: "Login", Password: "Password"},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockstorage(ctrl)
	uc := New(store)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidatePayload(tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestAuthUseCase_SignIn(t *testing.T) {
	type foundedUser struct {
		user *entities.UserWithPassword
		err  error
	}

	tests := []struct {
		name        string
		foundedUser foundedUser
		payload     entities.AuthPayload
		wantErr     bool
	}{
		{
			name: "error find request",
			foundedUser: foundedUser{
				user: nil,
				err:  errors.New("error find request"),
			},
			payload: entities.AuthPayload{
				Login: "login",
			},
			wantErr: true,
		},
		{
			name: "find return nil user",
			foundedUser: foundedUser{
				user: nil,
				err:  nil,
			},
			payload: entities.AuthPayload{
				Login: "login",
			},
			wantErr: true,
		},
		{
			name: "incorrect password",
			foundedUser: foundedUser{
				user: &entities.UserWithPassword{
					HashPassword: "incorrect password",
				},
				err: nil,
			},
			payload: entities.AuthPayload{
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "correct password",
			foundedUser: foundedUser{
				user: &entities.UserWithPassword{
					HashPassword: stubHashPassword,
				},
				err: nil,
			},
			payload: entities.AuthPayload{
				Login:    "login",
				Password: "password",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockstorage(ctrl)
			store.EXPECT().FindUserWithPassword(ctx, "login").Return(tt.foundedUser.user, tt.foundedUser.err)

			uc := New(store)

			got, err := uc.SignIn(ctx, tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func TestAuthUseCase_SignUp(t *testing.T) {
	type foundedUser struct {
		user *entities.UserWithPassword
		err  error
	}

	tests := []struct {
		name          string
		foundedUser   foundedUser
		payload       entities.AuthPayload
		wantErr       bool
		createUser    *entities.UserWithPassword
		createUserErr error
	}{
		{
			name: "find error req",
			foundedUser: foundedUser{
				err: errors.New("find error req"),
			},
			payload: entities.AuthPayload{
				Login: "login",
			},
			wantErr: true,
		},
		{
			name: "user is exist",
			foundedUser: foundedUser{
				user: &entities.UserWithPassword{},
			},
			payload: entities.AuthPayload{
				Login: "login",
			},
			wantErr: true,
		},
		{
			name:          "err create user req",
			payload:       entities.AuthPayload{Login: "login", Password: "password"},
			createUser:    &entities.UserWithPassword{Login: "login", HashPassword: stubHashPassword},
			createUserErr: errors.New("err create user req"),
			wantErr:       true,
		},
		{
			name:       "successful create user",
			payload:    entities.AuthPayload{Login: "login", Password: "password"},
			createUser: &entities.UserWithPassword{Login: "login", HashPassword: stubHashPassword},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockstorage(ctrl)
			store.EXPECT().FindUserWithPassword(ctx, "login").Return(tt.foundedUser.user, tt.foundedUser.err)

			if tt.createUser != nil {
				store.EXPECT().CreateUser(ctx, tt.createUser).Return(tt.createUserErr)
			}

			uc := New(store)

			got, err := uc.SignUp(ctx, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
