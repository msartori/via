package mock_guide_provider

import (
	"context"
	"via/internal/model"
)

type MockProvider struct {
	guides map[string]model.Guide
	err    error
}

func New(guides map[string]model.Guide, err error) *MockProvider {
	return &MockProvider{
		guides: guides,
		err:    err,
	}
}

func (m *MockProvider) GetGuide(ctx context.Context, id string) (model.Guide, error) {
	if m.err != nil {
		return model.Guide{}, m.err
	}
	return m.guides[id], nil
}
