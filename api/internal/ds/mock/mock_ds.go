package mock_ds

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDS struct {
	mock.Mock
}

func (m *MockDS) Set(ctx context.Context, key, value string, ttlSeconds int) error {
	args := m.Called(ctx, key, value, ttlSeconds)
	return args.Error(0)
}

func (m *MockDS) Get(ctx context.Context, key string) (bool, string, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockDS) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockDS) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}
