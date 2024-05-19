package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/authorization"
)

type API struct {
	service *services.AppService
}

func NewApi(s *services.AppService) *API {
	return &API{
		service: s,
	}
}

func (a *API) CreateRouter() (*mux.Router, error) {

	r := mux.NewRouter()

	envService, err := services.NewEnvService(".env")
	if err != nil {
		return r, err
	}

	jwtSecret := envService.Get("JWT_SECRET", "secret")
	jwtExpiration, _ := strconv.Atoi(envService.Get("JWT_EXPIRATION_HOURS", "24"))
	jwtService := authorization.NewJwtService(jwtSecret, jwtExpiration)

	authorizationService := authorization.NewAuthServie(jwtService)

	r.HandleFunc("/api/v1/users/", authorizationService.AdminMiddleware(a.AdminListHandler)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/", authorizationService.AdminMiddleware(a.AdminCreateHandler)).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/auth/admins/login/", a.LoginHandler).Methods(http.MethodPost)

	return r, nil
}
