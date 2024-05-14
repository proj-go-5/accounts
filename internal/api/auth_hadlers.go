package api

import (
	"encoding/json"
	"net/http"

	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/pkg/accountsio"
)

func (a *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.service.Auth.Login(loginRequest.Login, loginRequest.Password)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginResponse := dto.LoginResponse{Token: token}
	accountsio.MakeResponse(w, loginResponse, http.StatusOK)
}

func (a *API) TokenInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		accountsio.MakeResponse(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	jwt, err := a.service.Token.VerifyToken(token)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := a.service.Token.ExtractClaims(jwt)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	login := claims.Admin.Login

	cacheToken, exists, err := a.service.Cache.Get(login)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		accountsio.MakeResponse(w, "token expired", http.StatusUnauthorized)
		return
	}

	if cacheToken != token {
		accountsio.MakeResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accountsio.MakeResponse(w, claims, http.StatusOK)
}
