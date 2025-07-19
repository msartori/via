package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"
	mock_ds "via/internal/ds/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimiter_Allow(t *testing.T) {
	type fields struct {
		limit  int
		window time.Duration
	}
	type args struct {
		key       string
		incrCount int
		incrErr   error
		setErr    error
		wantAllow bool
		wantErr   bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "first request sets key",
			fields: fields{
				limit:  3,
				window: 5 * time.Second,
			},
			args: args{
				key:       "user1",
				incrCount: 1,
				wantAllow: true,
			},
		},
		{
			name: "multiple requests under limit",
			fields: fields{
				limit:  3,
				window: 5 * time.Second,
			},
			args: args{
				key:       "user2",
				incrCount: 2,
				wantAllow: true,
			},
		},
		{
			name: "limit exceeded",
			fields: fields{
				limit:  3,
				window: 5 * time.Second,
			},
			args: args{
				key:       "user3",
				incrCount: 4,
				wantAllow: false,
			},
		},
		{
			name: "error on Incr",
			fields: fields{
				limit:  3,
				window: 5 * time.Second,
			},
			args: args{
				key:     "user4",
				incrErr: errors.New("incr failed"),
				wantErr: true,
			},
		},
		{
			name: "error on Set when first request",
			fields: fields{
				limit:  3,
				window: 5 * time.Second,
			},
			args: args{
				key:       "user5",
				incrCount: 1,
				setErr:    errors.New("set failed"),
				wantErr:   true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDS := new(mock_ds.MockDS)

			if tt.args.incrErr != nil {
				mockDS.On("Incr", mock.Anything, tt.args.key).Return(0, tt.args.incrErr)
			} else {
				mockDS.On("Incr", mock.Anything, tt.args.key).Return(tt.args.incrCount, nil)
				if tt.args.incrCount == 1 {
					mockDS.On("Set", mock.Anything, tt.args.key, "1", int(tt.fields.window.Seconds())).
						Return(tt.args.setErr)
				}
			}

			rl := New("test", tt.fields.limit, tt.fields.window, mockDS)
			allowed, err := rl.Allow(context.Background(), tt.args.key)

			assert.Equal(t, tt.args.wantAllow, allowed)
			if tt.args.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDS.AssertExpectations(t)
		})
	}
}
