package middleware

/*
// Mocks
type MockLogger struct{}

func (m *MockLogger) Warn(ctx context.Context, args ...any)             {}
func (m *MockLogger) Error(ctx context.Context, err error, args ...any) {}
func (m *MockLogger) Info(ctx context.Context, args ...any)             {}
func (m *MockLogger) Debug(ctx context.Context, args ...any)            {}

type MockOIDCProvider struct{}

func TestJWTAuthMiddleware(t *testing.T) {
	testutil.InjectNoOpLogger()

	testutil.InjectMockJWTKey()
	// Mock verifyGoogleIDToken
	middleware.SetVerifyIDTokenFunc(func(ctx context.Context, idToken, clientID, issuer string) (bool, error) {
		if idToken == "valid-id-token" {
			return true, nil
		}
		return false, errors.New("invalid id token")
	})

	cfg := auth.OAuthConfig{
		ClientID:  "test-client",
		IDPIssuer: "https://issuer.example.com",
	}

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	t.Run("missing cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(cfg)(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, nextCalled)
	})

	t.Run("invalid jwt token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthTokenKey,
			Value: "invalid.jwt.token",
		})
		w := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(cfg)(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, nextCalled)
	})

	t.Run("invalid id token", func(t *testing.T) {
		token := createJWTToken(t, "invalid-id-token")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthTokenKey,
			Value: token,
		})
		w := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(cfg)(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, nextCalled)
	})

	t.Run("valid jwt and id token", func(t *testing.T) {
		nextCalled = false
		token := createJWTToken(t, "valid-id-token")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthTokenKey,
			Value: token,
		})
		w := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(cfg)(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, nextCalled)

		// Validate OperatorID in context
		opID := req.Context().Value(middleware.OperatorIDKey)
		assert.Equal(t, 42, opID)
	})
}

func createJWTToken(t *testing.T, idpIDToken string) string {
	t.Helper()
	claims := &middleware.Claims{
		OperatorID: 42,
		IDPIDToken: idpIDToken,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("dummy_secret"))
	assert.NoError(t, err)
	return signed
}
*/
