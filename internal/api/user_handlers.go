package api

import (
	"encoding/json"
	"net/http"

	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/proj-go-5/accounts/pkg/accountsio"
)

func (a *API) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId int64
	user, err := a.service.User.Save(&entities.AdminWithPassword{
		ID:       userId,
		Login:    createUserRequest.Login,
		Password: createUserRequest.Password,
	})

	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountsio.MakeResponse(w, user, http.StatusCreated)
}

func (a *API) UserListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.service.User.List()

	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accountsio.MakeResponse(w, users, http.StatusOK)
}
