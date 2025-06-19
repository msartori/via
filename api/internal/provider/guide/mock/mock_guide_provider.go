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
