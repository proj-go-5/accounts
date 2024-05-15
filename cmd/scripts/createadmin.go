package main

import (
	"log"
	"os"

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

	adminDbRepository, err := store.NewAdminDBRepository(envService)
	if err != nil {
		log.Println(err)
		return
	}
	defer adminDbRepository.Db.Close()

	adminService := services.NewAdminService(adminDbRepository)

	login := os.Args[1]
	password := os.Args[2]

	_, err = adminService.Save(&entities.AdminWithPassword{
		Login:    login,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}
}
