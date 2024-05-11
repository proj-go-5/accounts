package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/proj-go-5/accounts/internal/entities"
)

type UserDbRepository struct {
	db *sqlx.DB
}

func NewUserDBRepository(db *sqlx.DB) *UserDbRepository {
	return &UserDbRepository{db: db}
}

func (r *UserDbRepository) Save(a *entities.AdminWithPassword) (*entities.Admin, error) {

	var id int64

	q := fmt.Sprintf(`INSERT INTO admin (login, password) VALUES ('%v', '%v') RETURNING id`, a.Login, a.Password)

	err := r.db.QueryRow(q).Scan(&id)
	if err != nil {
		return nil, err
	}

	a.ID = id

	return a.WithoutPassword(), nil
}

func (r *UserDbRepository) List() ([]*entities.Admin, error) {
	var admins []*entities.Admin

	r.db.Select(&admins, "SELECT id, login FROM admin")
	return admins, nil
}

func (r *UserDbRepository) Get(login string) (*entities.AdminWithPassword, error) {
	var admins []*entities.AdminWithPassword

	q := fmt.Sprintf("SELECT id, login, password FROM admin WHERE login = '%v'", login)
	err := r.db.Select(&admins, q)

	if err != nil {
		return nil, err
	}

	if len(admins) == 0 {
		return nil, nil
	}

	if len(admins) > 1 {
		return nil, fmt.Errorf("'Get' returns multiple users with login %v", login)
	}

	return admins[0], nil
}
