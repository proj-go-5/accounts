package services

import (
	"errors"
	"fmt"

	"github.com/proj-go-5/accounts/internal/entities"
)

type UserRepository interface {
	Save(*entities.AdminWithPassword) (*entities.Admin, error)
	List() ([]*entities.Admin, error)
	Get(string) (*entities.AdminWithPassword, error)
}

type User struct {
	repository UserRepository
}

func NewUserService(r UserRepository) *User {
	return &User{repository: r}
}

func (u *User) List() ([]*entities.Admin, error) {
	return u.repository.List()
}

func (u *User) Save(user *entities.AdminWithPassword) (*entities.Admin, error) {
	login := user.Login

	dbUser, err := u.repository.Get(login)
	if err != nil {
		return nil, err
	}

	if dbUser != nil {
		return nil, errors.New(fmt.Sprintf("user with login %v already exists", login))
	}

	savedUser, err := u.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return savedUser, nil
}

func (u *User) Get(login string) (*entities.AdminWithPassword, error) {
	userWithPassword, err := u.repository.Get(login)
	if err != nil {
		return nil, err
	}
	return userWithPassword, nil
}
