package store

import (
	"github.com/proj-go-5/accounts/internal/entities"
	"sync"
)

type UserMemoryRepository struct {
	mx    sync.Mutex
	users []*entities.UserWithPassword
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{users: make([]*entities.UserWithPassword, 0)}
}

func (r *UserMemoryRepository) Save(u *entities.UserWithPassword) (*entities.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if u.ID == 0 {
		u.ID = int64(len(r.users) + 1)
	}

	r.users = append(r.users, u)
	return u.WithoutPassword(), nil
}

func (r *UserMemoryRepository) List() ([]*entities.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	users := make([]*entities.User, 0)

	for _, u := range r.users {
		users = append(users, u.WithoutPassword())
	}

	return users, nil
}

func (r *UserMemoryRepository) Get(login string) (*entities.UserWithPassword, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var dbUser *entities.UserWithPassword

	for _, user := range r.users {
		if user.Login == login {
			dbUser = user
			break
		}
	}
	return dbUser, nil
}
