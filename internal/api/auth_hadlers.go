package api

import (
	"encoding/json"
	"net/http"

	"github.com/proj-go-5/accounts/internal/api/dto"
)

func (a *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.service.Auth.Login(loginRequest.Login, loginRequest.Password)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginResponse := dto.LoginResponse{Token: token}
	a.makeResponse(w, loginResponse, http.StatusOK)
}

func (a *API) TokenInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		a.makeResponse(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	jwt, err := a.service.Token.VerifyToken(token)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := a.service.Token.ExtractClaims(jwt)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	login := claims.Login

	cacheToken, exists, err := a.service.Cache.Get(login)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		a.makeResponse(w, "token expired", http.StatusUnauthorized)
		return
	}

	if cacheToken != token {
		a.makeResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	a.makeResponse(w, claims, http.StatusOK)
}
