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
