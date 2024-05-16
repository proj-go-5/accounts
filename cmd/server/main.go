package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/proj-go-5/accounts/internal/api"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/authorization"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(NewEnvService),
		fx.Provide(NewDb),
		fx.Provide(store.NewAdminDBRepository),
		fx.Provide(NewRedisCli),
		fx.Provide(store.NewRedisCacheRepository),
		fx.Provide(services.NewCacheService),
		fx.Provide(services.NewHashService),
		fx.Provide(services.NewAdminService),
		fx.Provide(NewJwtService),
		fx.Provide(services.NewAuthService),
		fx.Provide(services.NewAppService),
		fx.Provide(api.New),

		fx.Invoke(runServer),
	)
}

func runServer(api *api.API, e *services.Env, cr services.CacheRepository, ar services.AdminRepository) {
	defer cr.Close()
	defer ar.Close()

	serverPort := e.Get("ACCOUNTS_SERVER_PORT", "8080")

	log.Printf("Runing servier on %v port\n", serverPort)

	r, err := api.CreateRouter()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
		log.Printf("Server run error: %s", err)
	}
}

func NewEnvService() (*services.Env, error) {
	envService, err := services.NewEnvService(".env")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return envService, nil
}

func NewJwtService(e *services.Env) (*authorization.JwtService, error) {
	jwtSecret := e.Get("JWT_SECRET", "secret")
	jwtExpiration, err := strconv.Atoi(e.Get("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, err
	}
	return authorization.NewJwtService(jwtSecret, jwtExpiration), nil
}

func NewDb(e *services.Env) (*sqlx.DB, error) {

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
	return db, nil
}

func NewRedisCli(e *services.Env) (*redis.Client, error) {
	redisAddres := fmt.Sprintf("%v:%v",
		e.Get("ACCOUNTS_REDIS_HOST", "localhost"),
		e.Get("ACCOUNTS_REDIS_PORT", "6379"),
	)

	redisDb, err := strconv.Atoi(e.Get("ACCOUNTS_REDIS_DB", "0"))
	if err != nil {
		return nil, err
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr:     redisAddres,
		Password: e.Get("ACCOUNTS_REDIS_PASSWORD", ""),
		DB:       redisDb,
	})

	err = redisCli.Ping().Err()
	if err != nil {
		return nil, err
	}
	return redisCli, nil
}
