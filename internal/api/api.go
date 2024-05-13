package api

import (
	"net/http"

	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/middlewares"

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

func (a *API) CreateRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users/", middlewares.Authorize(a.UserListHandler)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/", middlewares.Authorize(a.UserCreateHandler)).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/auth/admins/login/", a.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc(middlewares.VerifyTokenUrlPath, a.TokenInfoHandler).Methods(http.MethodGet)

	return r
}
