package services

import (
	"errors"
	"fmt"
	"github.com/proj-go-5/accounts/internal/entities"
)

type UserRepository interface {
	Save(*entities.UserWithPassword) (*entities.User, error)
	List() ([]*entities.User, error)
	Get(string) (*entities.UserWithPassword, error)
}

type User struct {
	repository UserRepository
}

func NewUserService(r UserRepository) *User {
	return &User{repository: r}
}

func (u *User) List() ([]*entities.User, error) {
	return u.repository.List()
}

func (u *User) Save(user *entities.UserWithPassword) (*entities.User, error) {
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

func (u *User) Get(login string) (*entities.UserWithPassword, error) {
	userWithPassword, err := u.repository.Get(login)
	if err != nil {
		return nil, err
	}
	return userWithPassword, nil
}
