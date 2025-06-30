package mock_guide_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockGuideProvider struct {
	mock.Mock
}

func (m *MockGuideProvider) GetGuideByViaGuideId(ctx context.Context, viaGuideId string) (model.Guide, error) {
	args := m.Called(ctx, viaGuideId)
	return args.Get(0).(model.Guide), args.Error(1)
}

func (m *MockGuideProvider) CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error) {
	args := m.Called(ctx, guide)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockGuideProvider) GetGuidesByStatus(ctx context.Context, status []string) ([]model.Guide, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]model.Guide), args.Error(1)
}

func (m *MockGuideProvider) UpdateGuide(ctx context.Context, guide model.Guide) error {
	args := m.Called(ctx, guide)
	return args.Error(0)
}

func (m *MockGuideProvider) GetGuideById(ctx context.Context, id int) (model.Guide, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Guide), args.Error(1)
}

func (m *MockGuideProvider) GetGuideHistory(ctx context.Context, guideId int) ([]model.GuideHistory, error) {
	args := m.Called(ctx, guideId)
	return args.Get(0).([]model.GuideHistory), args.Error(1)
}
