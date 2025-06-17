package mock_guide_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockGuideProvider struct {
	mock.Mock
}

func (m *MockGuideProvider) GetGuideByCode(ctx context.Context, code string) (model.GuideProcess, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(model.GuideProcess), args.Error(1)
}

func (m *MockGuideProvider) CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error) {
	args := m.Called(ctx, guide)
	return args.Get(0).(int), args.Error(1)
}
