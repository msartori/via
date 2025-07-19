package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"via/internal/auth"
	biz_operator "via/internal/biz/operator"
	"via/internal/ds"
	"via/internal/i18n"
	"via/internal/response"

	"via/internal/log"

	"golang.org/x/oauth2"
)

const (
	stateKey       = "state"
	redirectURIKey = "redirect_uri"
	codeKey        = "code"
)

func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, codeChallenge, challengeMethod, err := auth.Get().GenerateState(r.Context(), r.URL.Query().Get(redirectURIKey))
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to generate state")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		authURL := auth.Get().OAuth2Config.AuthCodeURL(state,
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
		oauth2Config := auth.Get().OAuth2Config
		authState, found, err := auth.Get().GetState(r.Context(), state)
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to get auth state")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
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

		signedToken, err := auth.GenerateAuthToken(operator, idToken, cfg)
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed to sign JWT token")
			response.WriteJSON(w, r, response.Response[any]{
				Message: i18n.Get(r, i18n.MsgInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		err = ds.Get().Set(r.Context(), idToken, token.AccessToken, cfg.JWTExpirationInSeconds)
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "unable to store access token", "error")
		}
		auth.SetAuthToken(w, signedToken, cfg)
		http.Redirect(w, r, authState.RedirectURI, http.StatusSeeOther)
	})
}

func LogOut(cfg auth.OAuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response.Response[any]{}
		token, err := auth.GetAuthToken(r)
		if err != nil {
			log.Get().Warn(r.Context(), "msg", "access token not found")
			resp.Message = i18n.Get(r, i18n.MsgAuthTokenNotFoud)
			response.WriteJSON(w, r, resp, http.StatusOK)
			return
		}

		auth.DelAuthToken(w, cfg)

		claims, _, err := auth.ParseTokenWithClaims(token)
		if err != nil {
			log.Get().Warn(r.Context(), "msg", "unable to parse JWT token", "error", err.Error())
			resp.Message = i18n.Get(r, i18n.MsgAuthTokenInvalid)
			response.WriteJSON(w, r, resp, http.StatusOK)
			return
		}
		ok, accessToken, err := ds.Get().Get(r.Context(), claims.IDPIDToken)
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "unable get access token from DS", "error")
			resp.Message = i18n.Get(r, i18n.MsgInternalServerError)
			response.WriteJSON(w, r, resp, http.StatusOK)
			return
		}
		if !ok {
			log.Get().Warn(r.Context(), "msg", "access token not found in DS")
			resp.Message = i18n.Get(r, i18n.MsgAccessTokenNotFound)
			response.WriteJSON(w, r, resp, http.StatusOK)
			return
		}
		response.WriteJSON(w, r, resp, http.StatusOK)

		go revokeToken(r.Context(), accessToken, cfg.RevokeTokenURL)

	})
}

var httpClient = http.DefaultClient

func revokeToken(ctx context.Context, token, uri string) {
	data := url.Values{}
	data.Set("token", token)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Get().Error(ctx, err, "msg", "unable to create revoke token request", "error", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "error in revoke token request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Get().Warn(ctx, "msg", "token revoke call error", "status", resp.StatusCode, "response", string(body))
	}
	log.Get().Info(ctx, "msg", "token revoke done")
}
