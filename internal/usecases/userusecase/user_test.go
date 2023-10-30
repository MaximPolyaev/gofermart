package userusecase

import (
	"context"
	"testing"

	"github.com/MaximPolyaev/gofermart/internal/usecases/userusecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetUserIDByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockstorage(ctrl)
	ctx := context.TODO()

	wantUserID := 1

	store.EXPECT().FindUserIDByLogin(ctx, "test").Return(wantUserID, nil)

	uc := New(store)
	got, err := uc.GetUserIDByLogin(ctx, "test")

	assert.NoError(t, err)
	assert.Equal(t, wantUserID, got)
}
