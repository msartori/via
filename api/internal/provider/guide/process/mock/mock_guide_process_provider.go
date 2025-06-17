package mock_guide_process_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockGuideProcessProvider struct {
	mock.Mock
}

func (m *MockGuideProcessProvider) GetGuideProcessByCode(ctx context.Context, code string) (model.GuideProcess, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(model.GuideProcess), args.Error(1)
}

func (m *MockGuideProcessProvider) CreateGuide(ctx context.Context, guide model.Guide) (int, error) {
	args := m.Called(ctx, guide)
	return args.Get(0).(int), args.Error(1)
}
