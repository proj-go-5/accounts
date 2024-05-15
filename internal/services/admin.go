package services

import (
	"errors"
	"fmt"

	"github.com/proj-go-5/accounts/internal/entities"
)

type AdminRepository interface {
	Save(*entities.AdminWithPassword) (*entities.Admin, error)
	List() ([]*entities.Admin, error)
	Get(string) (*entities.AdminWithPassword, error)
}

type Admin struct {
	repository AdminRepository
}

func NewAdminService(r AdminRepository) *Admin {
	return &Admin{repository: r}
}

func (a *Admin) List() ([]*entities.Admin, error) {
	return a.repository.List()
}

func (a *Admin) Save(admin *entities.AdminWithPassword) (*entities.Admin, error) {
	login := admin.Login

	dbAdmin, err := a.repository.Get(login)
	if err != nil {
		return nil, err
	}

	if dbAdmin != nil {
		return nil, errors.New(fmt.Sprintf("admin with login %v already exists", login))
	}

	savedAdmin, err := a.repository.Save(admin)
	if err != nil {
		return nil, err
	}

	return savedAdmin, nil
}

func (a *Admin) Get(login string) (*entities.AdminWithPassword, error) {
	adminWithPassword, err := a.repository.Get(login)
	if err != nil {
		return nil, err
	}
	return adminWithPassword, nil
}
