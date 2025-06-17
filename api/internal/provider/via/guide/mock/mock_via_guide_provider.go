package mock_via_guide_provider

import (
	"context"
	"via/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockViaGuideProvider struct {
	mock.Mock
}

func (m *MockViaGuideProvider) GetGuide(ctx context.Context, id string) (model.ViaGuide, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.ViaGuide), args.Error(1)
}
