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

func (m *MockDS) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockDS) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
