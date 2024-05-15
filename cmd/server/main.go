package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/proj-go-5/accounts/pkg/jwt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/proj-go-5/accounts/internal/api"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
)

func main() {
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
	defer db.Close()

	adminService := services.NewAdminService(store.NewAdminDBRepository(db))
	cacheService := services.NewCacheService(store.NewMemoryCacheRepository())

	jwtSecret := envService.Get("JWT_SECRET", "secret")
	jwtExpiration, _ := strconv.Atoi(envService.Get("JWT_EXPIRATION_HOURS", "24"))
	jwtService := jwt.NewJwtService(jwtSecret, jwtExpiration)

	appService := &services.AppService{
		Admin: adminService,
		Jwt:   jwtService,
		Cache: cacheService,
		Auth: services.NewAuthService(
			adminService, cacheService, jwtService,
		),
	}

	a := api.New(appService)

	r, err := a.CreateRouter()
	if err != nil {
		log.Fatal(err)
		return
	}

	serverPort := envService.Get("ACCOUNTS_SERVER_PORT", "8080")

	fmt.Printf("Runing servier on %v port\n", serverPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
		log.Printf("Server run error: %s", err)
	}
}
