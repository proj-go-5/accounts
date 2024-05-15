package services

import "github.com/proj-go-5/accounts/pkg/jwt"

type AppService struct {
	Admin *Admin
	Jwt   *jwt.Service
	Cache *Cache
	Auth  *Auth
}
