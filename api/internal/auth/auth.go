package auth

import (
	"context"
	"net/http"
	"sync"
	"time"
	"via/internal/ds"
	"via/internal/model"
	"via/internal/secret"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"via/internal/cache"

	jwt_key "via/internal/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

type Claims struct {
	OperatorID int    `json:"operatorId"`
	IDPIDToken string `json:"idpIdToken"`
	jwt.RegisteredClaims
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
		cfg.ClientSecret = secret.Get().Read(cfg.ClientSecretFile)
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
	RevokeTokenURL                string   `env:"REVOKE_TOKEN_URL" json:"revokeTokenUrl"`
	CookieSecure                  bool     `env:"COOKIE_SECURE" json:"secure"`
	CookieSameSiteNone            bool     `env:"COOKIE_SAME_SITE_NONE" json:"sameSiteNone"`
	CookieDomain                  string   `env:"COOKIE_DOMAIN" json:"domain"`
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

func GetAuthToken(r *http.Request) (token string, err error) {
	cookie, err := r.Cookie(AuthTokenKey)
	if err != nil {
		return
	}
	err = cookie.Valid()
	if err != nil {
		return
	}
	token = cookie.Value
	return
}

func writeAuthCookie(w http.ResponseWriter, token string, cfg OAuthConfig, expires time.Time) {
	sameSite := http.SameSiteLaxMode
	if cfg.CookieSameSiteNone {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     AuthTokenKey,
		Value:    token,
		Expires:  expires,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   cfg.CookieDomain,
	})
}

func SetAuthToken(w http.ResponseWriter, token string, cfg OAuthConfig) {
	writeAuthCookie(w, token, cfg, time.Now().Add(time.Duration(cfg.JWTExpirationInSeconds)*time.Second))
}

func DelAuthToken(w http.ResponseWriter, cfg OAuthConfig) {
	writeAuthCookie(w, "", cfg, time.Now().Add(-time.Hour))
}

func ParseTokenWithClaims(token string) (claims Claims, tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return jwt_key.GetPublicKey(), nil
	})
	return
}

func GenerateAuthToken(operator model.Operator, idToken string, cfg OAuthConfig) (string, error) {
	now := time.Now()
	claims := Claims{
		OperatorID: operator.ID,
		IDPIDToken: idToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWTClaimsIssuer,
			Audience:  jwt.ClaimStrings{cfg.JWTClaimsAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWTExpirationInSeconds) * time.Second)),
			ID:        uuid.New().String(),
			Subject:   operator.Account,
		},
	}
	tokenJWT := jwt.NewWithClaims(jwt_key.GetSigningMethod(), claims)
	return tokenJWT.SignedString(jwt_key.GetPrivateKey())
}
