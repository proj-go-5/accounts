package services

import (
	"errors"

	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/proj-go-5/accounts/pkg/authorization"
)

var defaultTtl = 60

type Auth struct {
	adminService *Admin
	cacheService *Cache
	jwtService   *authorization.JwtService
}

func NewAuthService(a *Admin, c *Cache, j *authorization.JwtService) *Auth {
	return &Auth{
		adminService: a,
		cacheService: c,
		jwtService:   j,
	}
}

func (a *Auth) CheckPassword(admin *entities.AdminWithPassword) (bool, error) {
	dbAdmin, err := a.adminService.Get(admin.Login)
	if err != nil {
		return false, err
	}

	if dbAdmin != nil {
		return dbAdmin.Password == admin.Password, nil
	}
	return false, nil
}

func (a *Auth) Login(login, password string) (string, error) {
	passwordOk, err := a.CheckPassword(&entities.AdminWithPassword{
		Login:    login,
		Password: password,
	})

	if err != nil {
		return "", err
	}

	if !passwordOk {
		return "", errors.New("wrong login or password")
	}

	token, exists, err := a.cacheService.Get(login)
	if err != nil {
		return "", err
	}

	if !exists {
		admin, err := a.adminService.Get(login)
		if err != nil {
			return "", err
		}

		adminClaims := &authorization.AdminClaims{
			ID:    admin.ID,
			Login: admin.Login,
		}
		token, err = a.jwtService.Generate(adminClaims)

		if err != nil {
			return "", err
		}

		if err = a.cacheService.Set(login, token, defaultTtl); err != nil {
			return "", err
		}
	}

	return token, nil
}
