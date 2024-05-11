package services

import (
	"errors"

	"github.com/proj-go-5/accounts/internal/entities"
)

var defaultTtl = 60

type Auth struct {
	userService  *User
	cacheService *Cache
	tokenService *Token
}

func NewAuthService(u *User, c *Cache, t *Token) *Auth {
	return &Auth{
		userService:  u,
		cacheService: c,
		tokenService: t,
	}
}

func (a *Auth) CheckPassword(user *entities.UserWithPassword) (bool, error) {
	dbUser, err := a.userService.Get(user.Login)
	if err != nil {
		return false, err
	}

	if dbUser != nil {
		return dbUser.Password == user.Password, nil
	}
	return false, nil
}

func (a *Auth) Login(login, password string) (string, error) {
	passwordOk, err := a.CheckPassword(&entities.UserWithPassword{
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
		user, err := a.userService.Get(login)
		if err != nil {
			return "", err
		}

		token, err = a.tokenService.Generate(user.WithoutPassword())

		if err != nil {
			return "", err
		}

		if err = a.cacheService.Set(login, token, defaultTtl); err != nil {
			return "", err
		}
	}

	return token, nil
}
