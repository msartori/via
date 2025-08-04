package mock_secret

import (
	"github.com/stretchr/testify/mock"
)

type MockSecret struct {
	mock.Mock
}

func (m *MockSecret) Read(path string) string {
	args := m.Called(path)
	return args.String(0)
}
