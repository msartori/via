package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte("super-secret-key")

type Claims struct {
	OperatorID string `json:"operatorId"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})
		if err != nil || !tkn.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "operatorId", claims.OperatorID))
		next.ServeHTTP(w, r)
	})
}
