package mock_guide_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockGuideProvider struct {
	mock.Mock
}

func (m *MockGuideProvider) GetGuide(ctx context.Context, id string) (model.Guide, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Guide), args.Error(1)
}
