package games

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ziscky/toggle-test/internal/sql"
	"github.com/ziscky/toggle-test/test/mocks"
)

func TestInitializeGameRequirements(t *testing.T) {
	p := &mocks.PersistInterface{}
	ctx := context.Background()

	tests := []struct {
		name    string
		mockFn  func()
		wantErr error
	}{
		{
			name: "success",
			mockFn: func() {
				p.On("CreateCards", mock.Anything, mock.Anything).Once().Return(
					nil)
			},
			wantErr: nil,
		},
		{
			name: "error creating cards",
			mockFn: func() {
				p.On("CreateCards", mock.Anything, mock.Anything).Once().Return(
					sql.ErrInternal)
			},
			wantErr: sql.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := InitializeGameRequirements(ctx, p)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
