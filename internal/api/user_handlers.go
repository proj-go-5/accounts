package api

import (
	"encoding/json"
	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/internal/entities"
	"net/http"
)

func (a *API) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId int64
	user, err := a.service.User.Save(&entities.UserWithPassword{
		ID:       userId,
		Login:    createUserRequest.Login,
		Password: createUserRequest.Password,
	})

	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.makeResponse(w, user, http.StatusCreated)
}

func (a *API) UserListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.service.User.List()

	if err != nil {
		a.makeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.makeResponse(w, users, http.StatusOK)
}
