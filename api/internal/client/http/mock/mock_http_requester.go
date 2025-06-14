package mock_http_requester

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockHttpRequester struct {
	mock.Mock
}

func (m *MockHttpRequester) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}
