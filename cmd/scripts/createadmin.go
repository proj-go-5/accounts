package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/proj-go-5/accounts/internal/entities"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Login and password parameters should be provided. Usage: ./createadmin.go <login> <password>")
	}

	envService, err := services.NewEnvService(".env")
	if err != nil {
		log.Println(err)
		return
	}

	dbDataSource := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v",
		envService.Get("ACCOUNTS_DB_USER", "accouunts"),
		envService.Get("ACCOUNTS_DB_PASSWORD", "accouunts"),
		envService.Get("ACCOUNTS_DB_NAME", "accouunts"),
		envService.Get("ACCOUNTS_DB_URL", "localhost"),
		envService.Get("ACCOUNTS_DB_PORT", "5432"),
		envService.Get("ACCOUNTS_DB_SSL_MODE", "disable"),
	)

	db, err := sqlx.Open("postgres", dbDataSource)
	if err != nil {
		log.Fatal(err)
		return
	}
	adminDbRepository := store.NewAdminDBRepository(db)
	defer db.Close()

	hashService := services.NewHashService()
	adminService := services.NewAdminService(adminDbRepository, hashService)

	login := os.Args[1]
	password := os.Args[2]

	_, err = adminService.Save(&entities.AdminWithPassword{
		Login:    login,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Admin '%v' successfully created\n", login)
}
