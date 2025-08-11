package mock_operator_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockOperatorProvider struct {
	mock.Mock
}

func (m *MockOperatorProvider) GetOperators(ctx context.Context) ([]model.Operator, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Operator), args.Error(1)
}

func (m *MockOperatorProvider) GetOperatorByAccount(ctx context.Context, account string) (model.Operator, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(model.Operator), args.Error(1)
}

func (m *MockOperatorProvider) GetOperatorById(ctx context.Context, id int) (model.Operator, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Operator), args.Error(1)
}
