package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildToken(t *testing.T) {
	token, err := BuildToken("login")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseClaims(t *testing.T) {
	type args struct {
		tokenString string
	}

	login := "login"
	validToken, err := BuildToken(login)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test case #1",
			args: args{
				tokenString: validToken,
			},
		},
		{
			name: "test case #2",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
					".eyJleHAiOjE2OTY4MzkzMzQsIlVzZXJMb2dpbiI6InRlc3QifQ" +
					".7_fW3fcHTWzMz93rXuqQTORNOXrru38Eb9tjnGLOuBE",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.args.tokenString)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
