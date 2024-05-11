package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/proj-go-5/accounts/internal/api"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
)

var defaultPort = "8080"

func main() {
	userService := services.NewUserService(store.NewUserMemoryRepository())
	cacheService := services.NewCacheService(store.NewMemoryCacheRepository())
	tokenService := services.NewTokenService()

	appService := &services.AppService{
		User:  userService,
		Token: tokenService,
		Cache: cacheService,
		Auth: services.NewAuthService(
			userService, cacheService, tokenService,
		),
	}

	a := api.New(appService)

	r := a.CreateRouter()

	port, ok := os.LookupEnv("ACCOUNTS_PORT")
	if !ok {
		fmt.Printf("'ACCOUNTS_PORT' env variable not found, runing the servier on a default port %s\n", defaultPort)
		port = defaultPort
	} else {
		fmt.Printf("Running the server on port %s\n", port)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Printf("Server run error: %s", err)
	}
}
