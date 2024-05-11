package api

import (
	"encoding/json"
	"net/http"

	"github.com/proj-go-5/accounts/internal/services"

	"github.com/gorilla/mux"
)

type API struct {
	service *services.AppService
}

func New(s *services.AppService) *API {
	return &API{
		service: s,
	}
}

func (a *API) makeResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *API) CreateRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/auth/admins/login/", a.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/admins/me/", a.TokenInfoHandler).Methods(http.MethodGet)

	return r
}
