package middleware

import (
	"context"
	"net/http"
	"via/internal/auth"
	"via/internal/i18n"
	jwt_key "via/internal/jwt"
	"via/internal/log"
	"via/internal/response"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
)

type OperatorIDKeyType string

const OperatorIDKey OperatorIDKeyType = "operatorId"

type Claims struct {
	OperatorID int    `json:"operatorId"`
	IDPIDToken string `json:"idpIdToken"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(cfg auth.OAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := response.Response[any]{
				Message: i18n.Get(r, i18n.MsgOperatorUnauthorized),
			}
			cookie, err := r.Cookie(auth.AuthTokenKey)
			if err != nil {
				log.Get().Warn(r.Context(), "msg", "missing auth token cookie", "error", err.Error())
				response.WriteJSON(w, r, res, http.StatusUnauthorized)
				return
			}
			claims := &Claims{}
			tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (any, error) {
				return jwt_key.GetPublicKey(), nil
			})
			if ok, err := verifyGoogleIDToken(r.Context(), claims.IDPIDToken, cfg.ClientID, cfg.IDPIssuer); err != nil || !ok {
				log.Get().Error(r.Context(), err, "msg", "failed to verify IDP token")
				response.WriteJSON(w, r, res, http.StatusUnauthorized)
				return
			}
			if err != nil {
				log.Get().Error(r.Context(), err, "msg", "failed to parse JWT token")
				response.WriteJSON(w, r, res, http.StatusUnauthorized)
				return
			}
			if !tkn.Valid {
				log.Get().Error(r.Context(), err, "msg", "invalid JWT token")
				response.WriteJSON(w, r, res, http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), OperatorIDKey, claims.OperatorID))
			next.ServeHTTP(w, r)
		})
	}
}

func verifyGoogleIDToken(ctx context.Context, idToken, clientID, issuer string) (bool, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed to create OIDC provider")
		return false, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	_, err = verifier.Verify(ctx, idToken)
	if err != nil {
		return false, err
	}
	return true, nil
}
