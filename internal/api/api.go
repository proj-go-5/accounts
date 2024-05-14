package api

import (
	"net/http"
	"strconv"

	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/env"
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

func (a *API) CreateRouter() (*mux.Router, error) {

	r := mux.NewRouter()

	envService, err := env.NewEnvService(".env")
	if err != nil {
		return r, err
	}

	jwtSecret := envService.Get("JWT_SECRET", "secret")
	jwtExpiration, _ := strconv.Atoi(envService.Get("JWT_EXPIRATION_HOURS", "24"))
	tokenService := services.NewTokenService(jwtSecret, jwtExpiration)

	authMiddlewareService := middlewares.NewAuthServie(tokenService)

	r.HandleFunc("/api/v1/users/", authMiddlewareService.Check(a.UserListHandler)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/", authMiddlewareService.Check(a.UserCreateHandler)).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/auth/admins/login/", a.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/admins/me/", a.TokenInfoHandler).Methods(http.MethodGet)

	return r, nil
}
