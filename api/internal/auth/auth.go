package auth

import (
	"sync"
	"time"
	"via/internal/secret"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

const AuthTokenKey string = "auth_token"

var (
	oauth2Config *oauth2.Config
	once         sync.Once
	// TODO move this to Redis or similar
	authCache = cache.New(5*time.Minute, 10*time.Minute)
)

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

func GetOAuth2Config(cfg OAuthConfig) *oauth2.Config {
	once.Do(func() {
		if cfg.ClientSecret == "" {
			cfg.ClientSecret = secret.ReadSecret(cfg.ClientSecretFile)
		}
		oauth2Config = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.AuthURL,
				TokenURL: cfg.TokenURL,
			},
		}
	})
	return oauth2Config
}

func generatePKCE() (codeVerifier, codeChallenge, method string, err error) {
	// Generate code_verifier (random string between 43-128 chars)
	b := make([]byte, 32)
	_, err = rand.Read(b)
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

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateState(redirectURI string) (string, string, string, error) {
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
	authCache.Set(state, authState, cache.DefaultExpiration)
	return state, codeChallenge, method, nil
}

func GetState(state string) (AuthState, bool) {
	if val, found := authCache.Get(state); found {
		if authState, ok := val.(AuthState); ok {
			return authState, true
		}
	}
	return AuthState{}, false
}
