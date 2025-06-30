package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"via/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	operatorEmailMap = map[string]string{ // simula DB de operadores
		"miguel.sartori@gmail.com": "1",
		"user2@gmail.com":          "2",
		"operador@gmail.com":       "3",
	}
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "https://via-api.fly.dev/auth/callback", // Cambiar por la URL real
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

func Login(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func LoginCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "failed to exchange token", http.StatusUnauthorized)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to get user info", http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "decode error", http.StatusInternalServerError)
		return
	}

	operatorID, ok := operatorEmailMap[userInfo.Email]
	if !ok {
		http.Error(w, "unauthorized user", http.StatusUnauthorized)
		return
	}

	claims := middleware.Claims{
		OperatorID: operatorID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := tokenJWT.SignedString(middleware.JWTKey)
	if err != nil {
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/operator", http.StatusFound)
}
