package store

import (
	"sync"

	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/proj-go-5/accounts/internal/services"
)

type AdminMemoryRepository struct {
	mx     sync.Mutex
	admins []*entities.AdminWithPassword
}

func NewAdminMemoryRepository() services.AdminRepository {
	return &AdminMemoryRepository{admins: make([]*entities.AdminWithPassword, 0)}
}

func (r *AdminMemoryRepository) Save(u *entities.AdminWithPassword) (*entities.Admin, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if u.ID == 0 {
		u.ID = int64(len(r.admins) + 1)
	}

	r.admins = append(r.admins, u)
	return u.WithoutPassword(), nil
}

func (r *AdminMemoryRepository) List() ([]*entities.Admin, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	admins := make([]*entities.Admin, 0)

	for _, u := range r.admins {
		admins = append(admins, u.WithoutPassword())
	}

	return admins, nil
}

func (r *AdminMemoryRepository) Get(login string) (*entities.AdminWithPassword, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var dbAdmin *entities.AdminWithPassword

	for _, admin := range r.admins {
		if admin.Login == login {
			dbAdmin = admin
			break
		}
	}
	return dbAdmin, nil
}
