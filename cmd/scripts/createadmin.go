package main

import (
	"log"

	"github.com/proj-go-5/accounts/internal/entities"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
)

func main() {
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

	_, err = adminService.Save(&entities.AdminWithPassword{
		Login:    "admin",
		Password: "admin",
	})

	if err != nil {
		log.Fatal(err)
	}
}
