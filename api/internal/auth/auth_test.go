package auth

import (
	"context"
	"errors"
	"testing"
	"via/internal/cache"
	mock_ds "via/internal/ds/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestAuthSingleton(t *testing.T) {
	old := Get()
	defer Set(old) // Restore after test

	// Set & Get
	Set(&Auth{})
	assert.NotNil(t, Get())
}

func TestNew(t *testing.T) {
	cfg := OAuthConfig{
		ClientID:     "client_id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"openid", "email"},
		AuthURL:      "https://auth.example.com",
		TokenURL:     "https://token.example.com",
	}
	expected := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.TokenURL,
		},
	}

	mockDS := new(mock_ds.MockDS)
	a := New(cfg, mockDS)

	assert.NotNil(t, a)
	assert.Equal(t, expected, a.OAuth2Config)
}

func TestGenerateStateAndGetState(t *testing.T) {
	ctx := context.Background()
	mockDS := new(mock_ds.MockDS)
	a := New(OAuthConfig{ClientSecret: "sec"}, mockDS)

	// Mock Set (simulate OK)
	mockDS.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	state, codeChallenge, method, err := a.GenerateState(ctx, "http://redirect")
	assert.NoError(t, err)
	assert.NotEmpty(t, state)
	assert.NotEmpty(t, codeChallenge)
	assert.Equal(t, "S256", method)

	mockDS.AssertCalled(t, "Set", mock.Anything, state, mock.Anything, mock.Anything)

	// Simulate fetching state successfully
	mockDS.On("Get", mock.Anything, state).
		Return(true, `{"CodeVerifier":"abc","RedirectURI":"http://redirect"}`, nil)

	authState, found, err := a.GetState(ctx, state)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "abc", authState.CodeVerifier)
	assert.Equal(t, "http://redirect", authState.RedirectURI)

	// Simulate cache miss
	mockDS.On("Get", mock.Anything, "missing").Return(false, "", nil)
	_, found, err = a.GetState(ctx, "missing")
	assert.NoError(t, err)
	assert.False(t, found)

	// Simulate cache error
	mockDS.On("Get", mock.Anything, "error").Return(false, "", errors.New("db error"))
	_, _, err = a.GetState(ctx, "error")
	assert.Error(t, err)
}

func TestGenerateStateErrors(t *testing.T) {
	ctx := context.Background()
	mockDS := new(mock_ds.MockDS)
	a := New(OAuthConfig{ClientSecret: "sec"}, mockDS)

	// Force Set() error
	mockDS.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("cache error"))

	_, _, _, err := a.GenerateState(ctx, "http://redirect")
	assert.Error(t, err)
}

func TestGenerateState_Errors(t *testing.T) {
	ctx := context.Background()

	// Mock DS & Cache
	mockDS := new(mock_ds.MockDS)
	c := cache.New(mockDS)

	// Instance under test
	a := &Auth{
		cache: c,
	}

	originalRandRead := randReadState

	randReadState = func(b []byte) (int, error) {
		return 0, errors.New("state generation generateState() failed")
	}
	_, _, _, err := a.GenerateState(ctx, "https://example.com/callback")
	assert.EqualError(t, err, "state generation generateState() failed")

	randReadState = originalRandRead

	originalRandRead = randReadPKCE

	randReadPKCE = func(b []byte) (int, error) {
		return 0, errors.New("state generation generatePKCE() failed")
	}
	_, _, _, err = a.GenerateState(ctx, "https://example.com/callback")
	assert.EqualError(t, err, "state generation generatePKCE() failed")

	randReadPKCE = originalRandRead
}
