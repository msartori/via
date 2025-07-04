package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"via/internal/auth"
	biz_operator "via/internal/biz/operator"
	"via/internal/i18n"
	jwt_key "via/internal/jwt"
	"via/internal/middleware"
	"via/internal/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"via/internal/log"
)

const (
	stateKey       = "state"
	redirectURIKey = "redirect_uri"
	codeKey        = "code"
)

func Login(cfg auth.OAuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, codeChallenge, challengeMethod, err := auth.GenerateState(r.URL.Query().Get(redirectURIKey))
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to generate state")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		authURL := auth.GetOAuth2Config(cfg).AuthCodeURL(state,
			oauth2.AccessTypeOffline,
			oauth2.SetAuthURLParam("code_challenge_method", challengeMethod),
			oauth2.SetAuthURLParam("code_challenge", codeChallenge))
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	})
}

func LoginCallback(cfg auth.OAuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get(codeKey)
		state := r.URL.Query().Get(stateKey)
		oauth2Config := auth.GetOAuth2Config(cfg)
		authState, found := auth.GetState(state)
		if !found {
			log.Get().Warn(r.Context(), "msg", "auth state not found")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgAuthStateNotFound),
			}, http.StatusBadRequest)
			return
		}
		token, err := oauth2Config.Exchange(r.Context(),
			code,
			oauth2.SetAuthURLParam("code_verifier", authState.CodeVerifier))

		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to exchange token")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgAuthFailedToExchangeToken),
			}, http.StatusUnauthorized)
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			log.Get().Error(r.Context(), nil, "msg", "missing id_token in token response")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgAuthFailedToExchangeToken),
			}, http.StatusUnauthorized)
			return
		}

		client := oauth2Config.Client(r.Context(), token)
		resp, err := client.Get(cfg.UserInfoURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Get().Error(r.Context(), err, "msg", "failed to get user info", "resp_status", resp.StatusCode)
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgAuthFailedToGetUserInfo),
			}, http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()
		var userInfo struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to decode user info response")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		operator, err := biz_operator.GetOperatorByAccount(r.Context(), userInfo.Email)
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to get operator by account")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		if operator.ID == 0 {
			log.Get().Error(r.Context(), err, "msg", "operator not found", "account", userInfo.Email)
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgOperatorInvalid),
			}, http.StatusUnauthorized)
			return
		}
		if !operator.Enabled {
			log.Get().Error(r.Context(), err, "msg", "operator is not active", "account", userInfo.Email)
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgOperatorUnauthorized),
			}, http.StatusUnauthorized)
			return
		}
		now := time.Now()
		claims := middleware.Claims{
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
		signedToken, err := tokenJWT.SignedString(jwt_key.GetPrivateKey())
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to sign JWT token")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     auth.AuthTokenKey,
			Value:    signedToken,
			Expires:  time.Now().Add(time.Duration(cfg.JWTRefreshExpirationInSeconds) * time.Second),
			HttpOnly: true,
			Secure:   cfg.SecureCookie,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})
		http.Redirect(w, r, authState.RedirectURI, http.StatusSeeOther)
	})
}
