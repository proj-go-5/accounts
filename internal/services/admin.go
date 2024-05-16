package services

import (
	"fmt"

	"github.com/proj-go-5/accounts/internal/entities"
)

type AdminRepository interface {
	Save(*entities.AdminWithPassword) (*entities.Admin, error)
	List() ([]*entities.Admin, error)
	Get(string) (*entities.AdminWithPassword, error)
	Close()
}

type Admin struct {
	repository  AdminRepository
	hashService *Hash
}

func NewAdminService(r AdminRepository, h *Hash) *Admin {
	return &Admin{repository: r, hashService: h}
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
		return nil, fmt.Errorf("admin with login %v already exists", login)
	}

	hash, err := a.hashService.HashPassword(admin.Password)
	if err != nil {
		return nil, err
	}

	admin.Password = hash
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
