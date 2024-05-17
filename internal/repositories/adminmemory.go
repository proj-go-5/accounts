package store

import (
	"sync"

	"github.com/proj-go-5/accounts/internal/entities"
)

type UserMemoryRepository struct {
	mx    sync.Mutex
	users []*entities.AdminWithPassword
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{users: make([]*entities.AdminWithPassword, 0)}
}

func (r *UserMemoryRepository) Save(u *entities.AdminWithPassword) (*entities.Admin, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if u.ID == 0 {
		u.ID = int64(len(r.users) + 1)
	}

	r.users = append(r.users, u)
	return u.WithoutPassword(), nil
}

func (r *UserMemoryRepository) List() ([]*entities.Admin, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	users := make([]*entities.Admin, 0)

	for _, u := range r.users {
		users = append(users, u.WithoutPassword())
	}

	return users, nil
}

func (r *UserMemoryRepository) Get(login string) (*entities.AdminWithPassword, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var dbUser *entities.AdminWithPassword

	for _, user := range r.users {
		if user.Login == login {
			dbUser = user
			break
		}
	}
	return dbUser, nil
}
