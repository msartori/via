package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuth2Cfg interface {
	Client(ctx context.Context, t *oauth2.Token) *http.Client
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
}
