package ordersusecase

import (
	"testing"

	"github.com/MaximPolyaev/gofermart/internal/usecases/ordersusecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrdersUseCase_ValidateNumber(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		wantErr bool
	}{
		{
			name:    "empty number",
			number:  "",
			wantErr: true,
		},
		{
			name:    "valid number #1",
			number:  "123467890901",
			wantErr: false,
		},
		{
			name:    "valid number #2",
			number:  "003467890905",
			wantErr: false,
		},
		{
			name:    "no valid number #1",
			number:  "123467890902",
			wantErr: true,
		},
		{
			name:    "no valid number #2",
			number:  "003467890904",
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockstorage(ctrl)
	uc := New(store)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateNumber(tt.number)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
