package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/proj-go-5/accounts/pkg/accountsio"
)

func (a *API) AdminCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId int64
	admin, err := a.service.Admin.Save(&entities.AdminWithPassword{
		ID:       userId,
		Login:    createUserRequest.Login,
		Password: createUserRequest.Password,
	})

	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountsio.MakeResponse(w, admin, http.StatusCreated)
}

func (a *API) AdminListHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("UserId").(string)
	userLogin, _ := r.Context().Value("UserLogin").(string)

	log.Printf("%v for Admin id = %v username = %v", r.URL, userId, userLogin)

	admin, err := a.service.Admin.List()

	if err != nil {
		accountsio.MakeResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accountsio.MakeResponse(w, admin, http.StatusOK)
}
