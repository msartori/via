package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"via/internal/auth"
	biz_operator "via/internal/biz/operator"
	mock_ds "via/internal/ds/mock"
	"via/internal/i18n"
	jwt_key "via/internal/jwt"
	jwt_key_mock "via/internal/jwt/mock"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/model"
	"via/internal/testutil"

	"via/internal/response"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auth_mock "via/internal/auth/mock"

	operator_provider "via/internal/provider/operator"
	mock_operator_provider "via/internal/provider/operator/mock"
)

func TestLogin_Success(t *testing.T) {
	// Arrange: mock DS returns expected state without error.
	mockDS := new(mock_ds.MockDS)
	mockDS.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	a := auth.New(auth.OAuthConfig{
		RedirectURL:  "https://example.com/callback",
		AuthURL:      "https://auth.example.com/authorize",
		ClientSecret: "mySecret",
	},
		mockDS)

	auth.Set(a)

	req := httptest.NewRequest(http.MethodGet, "/login?redirect_uri=https://client.app/callback", nil)
	w := httptest.NewRecorder()

	// Act
	Login().ServeHTTP(w, req)

	// Assert
	res := w.Result()
	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	location := res.Header.Get("Location")
	assert.Contains(t, location, "https://auth.example.com/authorize")
	assert.Contains(t, location, "code_challenge")
	assert.Contains(t, location, "code_challenge_method")
	assert.Contains(t, location, "state")

	mockDS.AssertExpectations(t)
}

func TestLogin_GenerateStateError(t *testing.T) {
	mockDS := new(mock_ds.MockDS)
	mockDS.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("DB error"))

	a := auth.New(auth.OAuthConfig{
		RedirectURL:  "https://example.com/callback",
		AuthURL:      "https://auth.example.com/authorize",
		ClientSecret: "mySecret",
	},
		mockDS)

	auth.Set(a)

	req := httptest.NewRequest(http.MethodGet, "/login?redirect_uri=https://client.app/callback", nil)
	w := httptest.NewRecorder()

	log.Set(&mock_log.MockNoOpLogger{})

	Login().ServeHTTP(w, req)

	// Assert
	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var resp response.Response[any]
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Message)
	mockDS.AssertExpectations(t)
}

func TestLoginCallback_Success(t *testing.T) {
	// Mock DS to return valid state
	t.Cleanup(func() {
		jwt_key.Reset()
		biz_operator.ClearCache()
	})
	mockDS := new(mock_ds.MockDS)
	mockDS.On("Get", mock.Anything, mock.Anything).
		Return(true, `{"CodeVerifier":"test_verifier","RedirectURI":"https://client.app/callback"}`, nil)

	// Mock OAuth2 config behavior
	mockOAuthCfg := new(auth_mock.MockOAuth2Cfg)
	exchangedToken := &oauth2.Token{}
	exchangedToken = exchangedToken.WithExtra(map[string]interface{}{
		"id_token": "dummy_id_token",
	})
	mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
		Return(exchangedToken, nil)

	mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
		Return(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) *http.Response {
				userInfo := map[string]string{"email": "test@example.com"}
				body, _ := json.Marshal(userInfo)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}
			}),
		})

	// Inject auth with mocks
	a := auth.New(auth.OAuthConfig{ClientSecret: "mySecret"}, mockDS)
	a.OAuth2Config = mockOAuthCfg
	auth.Set(a)

	// Inject mock logger
	log.Set(&mock_log.MockNoOpLogger{})

	// Mock Operator Provider
	mockOperatorProvider := new(mock_operator_provider.MockOperatorProvider)
	mockOperatorProvider.On("GetOperatorByAccount", mock.Anything, "test@example.com").
		Return(model.Operator{
			ID:      1,
			Account: "test@example.com",
			Enabled: true,
		}, nil)
	operator_provider.Set(mockOperatorProvider)

	req := httptest.NewRequest(http.MethodGet, "/callback?code=test_code&state=test_state", nil)
	w := httptest.NewRecorder()
	testutil.InjectMockJWTKey()
	LoginCallback(auth.OAuthConfig{
		JWTClaimsIssuer:               "test-issuer",
		JWTClaimsAudience:             "test-audience",
		JWTExpirationInSeconds:        3600,
		JWTRefreshExpirationInSeconds: 7200,
		SecureCookie:                  false,
		UserInfoURL:                   "https://userinfo.test",
	}).ServeHTTP(w, req)

	// Assert
	res := w.Result()
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.NotEmpty(t, res.Header.Get("Set-Cookie"))
	assert.Equal(t, "https://client.app/callback", res.Header.Get("Location"))

	mockDS.AssertExpectations(t)
	mockOAuthCfg.AssertExpectations(t)
	mockOperatorProvider.AssertExpectations(t)
}

// Helper to mock HTTP Client RoundTrip
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestLoginCallback_Errors(t *testing.T) {
	type testCase struct {
		name           string
		setupMocks     func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, mockOperatorProvider *mock_operator_provider.MockOperatorProvider)
		expectedStatus int
		expectedMsg    string
		customJWTKey   func()
	}

	tests := []testCase{
		{
			name: "error getting auth state",
			setupMocks: func(mockDS *mock_ds.MockDS, _ *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(false, "", errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    i18n.MsgInternalServerError,
		},
		{
			name: "auth state not found",
			setupMocks: func(mockDS *mock_ds.MockDS, _ *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(false, "", nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    i18n.MsgAuthStateNotFound,
		},
		{
			name: "error exchanging token",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(&oauth2.Token{}, errors.New("exchange error"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    i18n.MsgAuthFailedToExchangeToken,
		},
		{
			name: "missing id_token",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(&oauth2.Token{}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    i18n.MsgAuthFailedToExchangeToken,
		},
		{
			name: "error getting user info",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]interface{}{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}
						}),
					})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    i18n.MsgAuthFailedToGetUserInfo,
		},
		{
			name: "error decoding user info",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, _ *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]interface{}{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("invalid json"))}
						}),
					})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    i18n.MsgInternalServerError,
		},
		{
			name: "error getting operator by account",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, mockOperatorProvider *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]interface{}{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							userInfo := map[string]string{"email": "test@example.com"}
							body, _ := json.Marshal(userInfo)
							return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(body))}
						}),
					})
				mockOperatorProvider.On("GetOperatorByAccount", mock.Anything, "test@example.com").
					Return(model.Operator{}, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    i18n.MsgInternalServerError,
		},
		{
			name: "operator not found",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, mockOperatorProvider *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]any{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							userInfo := map[string]string{"email": "test@example.com"}
							body, _ := json.Marshal(userInfo)
							return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(body))}
						}),
					})
				mockOperatorProvider.On("GetOperatorByAccount", mock.Anything, "test@example.com").
					Return(model.Operator{ID: 0}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    i18n.MsgOperatorInvalid,
		},
		{
			name: "operator not enabled",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, mockOperatorProvider *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)
				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]any{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							userInfo := map[string]string{"email": "test@example.com"}
							body, _ := json.Marshal(userInfo)
							return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(body))}
						}),
					})
				mockOperatorProvider.On("GetOperatorByAccount", mock.Anything, "test@example.com").
					Return(model.Operator{ID: 1, Enabled: false, Account: "test@example.com"}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    i18n.MsgOperatorUnauthorized,
		}, {
			name: "error signing jwt token",
			setupMocks: func(mockDS *mock_ds.MockDS, mockOAuthCfg *auth_mock.MockOAuth2Cfg, mockOperatorProvider *mock_operator_provider.MockOperatorProvider) {
				mockDS.On("Get", mock.Anything, mock.Anything).
					Return(true, `{"CodeVerifier":"verifier","RedirectURI":"https://redirect"}`, nil)

				exchangedToken := &oauth2.Token{}
				exchangedToken = exchangedToken.WithExtra(map[string]any{"id_token": "token"})
				mockOAuthCfg.On("Exchange", mock.Anything, "test_code", mock.Anything).
					Return(exchangedToken, nil)
				mockOAuthCfg.On("Client", mock.Anything, exchangedToken).
					Return(&http.Client{
						Transport: roundTripFunc(func(*http.Request) *http.Response {
							userInfo := map[string]string{"email": "test@example.com"}
							body, _ := json.Marshal(userInfo)
							return &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(bytes.NewBuffer(body)),
							}
						}),
					})
				mockOperatorProvider.On("GetOperatorByAccount", mock.Anything, "test@example.com").
					Return(model.Operator{ID: 1, Enabled: true, Account: "test@example.com"}, nil)
			},
			customJWTKey: func() {
				// Inject an invalid private key to force signing error
				jwt_key.Reset()
				jwt_key.Init(jwt_key.JWTConfig{
					PrivateKey:    jwt_key_mock.GetPrivateKey(),
					PublicKey:     jwt_key_mock.GetPublicKey(),
					SigningMethod: &jwt_key_mock.MockSigner{},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    i18n.MsgInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDS := new(mock_ds.MockDS)
			mockOAuthCfg := new(auth_mock.MockOAuth2Cfg)
			mockOperatorProvider := new(mock_operator_provider.MockOperatorProvider)

			tt.setupMocks(mockDS, mockOAuthCfg, mockOperatorProvider)

			a := auth.New(auth.OAuthConfig{ClientSecret: "mySecret"}, mockDS)
			a.OAuth2Config = mockOAuthCfg
			auth.Set(a)
			operator_provider.Set(mockOperatorProvider)
			log.Set(&mock_log.MockNoOpLogger{})

			if tt.customJWTKey != nil {
				tt.customJWTKey()
			}

			req := httptest.NewRequest(http.MethodGet, "/callback?code=test_code&state=test_state", nil)
			w := httptest.NewRecorder()
			LoginCallback(auth.OAuthConfig{
				UserInfoURL: "https://userinfo.test",
			}).ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			var resp response.Response[any]
			err := json.NewDecoder(res.Body).Decode(&resp)
			assert.NoError(t, err)
			assert.Equal(t, i18n.Get(req, tt.expectedMsg), resp.Message)
			t.Cleanup(func() {
				jwt_key.Reset()
				biz_operator.ClearCache()
			})
		})
	}
}
