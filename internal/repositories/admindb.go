package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/proj-go-5/accounts/internal/services"
)

type AdminDbRepository struct {
	Db *sqlx.DB
}

func NewAdminDBRepository(e *services.Env) (*AdminDbRepository, error) {
	dbDataSource := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v",
		e.Get("ACCOUNTS_DB_USER", "accouunts"),
		e.Get("ACCOUNTS_DB_PASSWORD", "accouunts"),
		e.Get("ACCOUNTS_DB_NAME", "accouunts"),
		e.Get("ACCOUNTS_DB_URL", "localhost"),
		e.Get("ACCOUNTS_DB_PORT", "5432"),
		e.Get("ACCOUNTS_DB_SSL_MODE", "disable"),
	)

	db, err := sqlx.Open("postgres", dbDataSource)
	if err != nil {
		return nil, err
	}
	return &AdminDbRepository{Db: db}, nil
}

func (r *AdminDbRepository) Save(a *entities.AdminWithPassword) (*entities.Admin, error) {

	var id int64

	r.Db.QueryRow("INSERT INTO admin (login, password) VALUES ($1, $2) RETURNING id",
		a.Login, a.Password).Scan(&id)

	a.ID = id

	return a.WithoutPassword(), nil
}

func (r *AdminDbRepository) List() ([]*entities.Admin, error) {
	admins := make([]*entities.Admin, 0)

	r.Db.Select(&admins, "SELECT id, login FROM admin")
	return admins, nil
}

func (r *AdminDbRepository) Get(login string) (*entities.AdminWithPassword, error) {
	var admins []entities.AdminWithPassword

	err := r.Db.Select(&admins, "SELECT id, login, password FROM admin WHERE login = $1", login)

	if err != nil {
		return nil, err
	}

	if len(admins) == 0 {
		return nil, nil
	}

	if len(admins) > 1 {
		return nil, fmt.Errorf("'Get' returns multiple users with login %v", login)
	}

	return &admins[0], nil
}
