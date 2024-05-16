package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/proj-go-5/accounts/internal/api"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/authorization"
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

	cacheStore, err := store.NewRedisCacheRepository(envService)
	if err != nil {
		log.Println(err)
		return
	}
	defer cacheStore.Cli.Close()

	hashService := services.NewHashService()
	adminService := services.NewAdminService(adminDbRepository, hashService)
	cacheService := services.NewCacheService(cacheStore)

	jwtSecret := envService.Get("JWT_SECRET", "secret")
	jwtExpiration, _ := strconv.Atoi(envService.Get("JWT_EXPIRATION_HOURS", "24"))
	jwtService := authorization.NewJwtService(jwtSecret, jwtExpiration)

	appService := &services.AppService{
		Admin: adminService,
		Jwt:   jwtService,
		Cache: cacheService,
		Auth: services.NewAuthService(
			adminService, cacheService, jwtService, hashService,
		),
	}

	a := api.New(appService)

	r, err := a.CreateRouter()
	if err != nil {
		log.Fatal(err)
		return
	}

	serverPort := envService.Get("ACCOUNTS_SERVER_PORT", "8080")

	log.Printf("Runing servier on %v port\n", serverPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
		log.Printf("Server run error: %s", err)
	}
}
