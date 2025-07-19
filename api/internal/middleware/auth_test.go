package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/auth"
	"via/internal/model"
	"via/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	testutil.InjectNoOpLogger()
	testutil.InjectMockJWTKey()
	type testCase struct {
		name                 string
		token                string
		verifyFunc           func(ctx context.Context, idToken, clientID, issuer string) (bool, error)
		expectedStatus       int
		expectNextHandlerRun bool
	}

	cfg := auth.OAuthConfig{
		ClientID:  "client-id",
		IDPIssuer: "https://issuer",
	}

	validToken, _ := auth.GenerateAuthToken(model.Operator{ID: 5}, "idp-token", auth.OAuthConfig{JWTExpirationInSeconds: 10})

	cases := []testCase{
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token parse",
			token:          "invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "id token verify fails",
			token:          validToken,
			expectedStatus: http.StatusUnauthorized,
			verifyFunc: func(ctx context.Context, idToken string, clientID string, issuer string) (bool, error) {
				return false, errors.New("verify error")
			},
		},
		{
			name:           "success",
			token:          validToken,
			expectedStatus: http.StatusOK,
			verifyFunc: func(ctx context.Context, idToken string, clientID string, issuer string) (bool, error) {
				return true, nil
			},
			expectNextHandlerRun: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.verifyFunc != nil {
				verifiIdTokenFunc = tc.verifyFunc
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.token != "" {
				req.AddCookie(&http.Cookie{
					Name:  auth.AuthTokenKey,
					Value: tc.token,
				})
			}
			rr := httptest.NewRecorder()

			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
				val := r.Context().Value(OperatorIDKey)
				if tc.expectNextHandlerRun {
					assert.Equal(t, 5, val)
				}
			})

			handler := Auth(cfg)(next)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectNextHandlerRun, called)
		})
	}
}
