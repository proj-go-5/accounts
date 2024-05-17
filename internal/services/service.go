package services

import "github.com/proj-go-5/accounts/pkg/authorization"

type AppService struct {
	Admin *Admin
	Jwt   *authorization.JwtService
	Cache *Cache
	Auth  *Auth
}

func NewAppService(a *Admin, j *authorization.JwtService, c *Cache, auth *Auth) *AppService {
	return &AppService{
		Admin: a,
		Jwt:   j,
		Cache: c,
		Auth:  auth,
	}
}
