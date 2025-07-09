package auth

import (
	"context"
	"sync"
	"via/internal/ds"
	"via/internal/secret"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"via/internal/cache"

	"golang.org/x/oauth2"
)

const AuthTokenKey string = "auth_token"
const cacheTTL = 5 * 60

var (
	instance *Auth
	mutex    = &sync.RWMutex{}
)

type Auth struct {
	OAuth2Config OAuth2Cfg
	cache        *cache.Cache // Use a suitable cache implementation, e.g., Redis
}

func Get() *Auth {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(auth *Auth) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = auth
}

func New(cfg OAuthConfig, ds ds.DS) *Auth {
	if cfg.ClientSecret == "" {
		cfg.ClientSecret = secret.ReadSecret(cfg.ClientSecretFile)
	}
	return &Auth{
		cache: cache.New(ds),
		OAuth2Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.AuthURL,
				TokenURL: cfg.TokenURL,
			},
		},
	}
}

// AuthState holds the state for OAuth2 authentication, including the code verifier and redirect URI.
type AuthState struct {
	CodeVerifier string
	RedirectURI  string
}

type OAuthConfig struct {
	ClientID                      string   `env:"CLIENT_ID" json:"-"`
	ClientSecret                  string   `env:"CLIENT_SECRET" json:"-"`
	ClientSecretFile              string   `env:"CLIENT_SECRET_FILE" json:"clientSecretFile"`
	RedirectURL                   string   `env:"REDIRECT_URL" json:"redirectUrl"`
	Scopes                        []string `env:"SCOPES" json:"scopes"`
	AuthURL                       string   `env:"AUTH_URL" json:"authUrl"`
	TokenURL                      string   `env:"TOKEN_URL" json:"tokenUrl"`
	UserInfoURL                   string   `env:"USER_INFO_URL" json:"userInfoUrl"`
	SecureCookie                  bool     `env:"SECURE_COOKIE" json:"secure"`
	JWTClaimsIssuer               string   `env:"JWT_CLAIMS_ISSUER" json:"jwtClaimsIssuer"`
	JWTClaimsAudience             string   `env:"JWT_CLAIMS_AUDIENCE" json:"jwtClaimsAudience"`
	JWTExpirationInSeconds        int      `env:"JWT_EXPIRATION_IN_SECONDS" json:"jwtExpirationInSeconds"`
	JWTRefreshExpirationInSeconds int      `env:"JWT_REFRESH_EXPIRATION_IN_SECONDS" json:"jwtRefreshExpirationInSeconds"`
	IDPIssuer                     string   `env:"IDP_ISSUER" json:"idpIssuer"`
}

var randReadPKCE = rand.Read

func generatePKCE() (codeVerifier, codeChallenge, method string, err error) {
	// Generate code_verifier (random string between 43-128 chars)
	b := make([]byte, 32)
	_, err = randReadPKCE(b)
	if err != nil {
		return "", "", "", err
	}
	codeVerifier = base64.RawURLEncoding.EncodeToString(b)

	// Generate code_challenge (SHA256 of code_verifier, base64-url-encoded)
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)
	codeChallenge = base64.RawURLEncoding.EncodeToString(hashed)

	return codeVerifier, codeChallenge, "S256", nil
}

var randReadState = rand.Read

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := randReadState(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (a *Auth) GenerateState(ctx context.Context, redirectURI string) (string, string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", "", err
	}
	codeVerifier, codeChallenge, method, err := generatePKCE()
	if err != nil {
		return "", "", "", err
	}
	authState := AuthState{
		CodeVerifier: codeVerifier,
		RedirectURI:  redirectURI}

	err = a.cache.Set(ctx, state, authState, cacheTTL)
	if err != nil {
		return "", "", "", err
	}
	return state, codeChallenge, method, nil
}

func (a *Auth) GetState(ctx context.Context, state string) (AuthState, bool, error) {
	var authState AuthState
	found, err := a.cache.Get(ctx, state, &authState)
	if err != nil {
		return AuthState{}, false, err
	}
	// Check if the value is found and is of type AuthState
	if !found {
		return AuthState{}, false, nil // Not found
	}
	return authState, true, nil
}
