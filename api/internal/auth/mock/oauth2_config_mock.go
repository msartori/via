package auth_mock

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockOAuth2Cfg struct {
	mock.Mock
}

func (m *MockOAuth2Cfg) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	args := m.Called(ctx, token)
	return args.Get(0).(*http.Client)
}

func (m *MockOAuth2Cfg) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	args := m.Called(ctx, code, opts)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockOAuth2Cfg) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	args := m.Called(state, opts)
	return args.String(0)
}
