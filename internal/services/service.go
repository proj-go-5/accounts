package services

import "github.com/proj-go-5/accounts/pkg/authorization"

type AppService struct {
	Admin *Admin
	Jwt   *authorization.JwtService
	Cache *Cache
	Auth  *Auth
}
