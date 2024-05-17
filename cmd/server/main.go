package main

import (
	"context"
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
	).Run()
}

func runServer(lifecycle fx.Lifecycle, api *api.API, e *services.Env) {
	serverPort := e.Get("ACCOUNTS_SERVER_PORT", "8080")

	log.Printf("Runing servier on %v port\n", serverPort)

	r, err := api.CreateRouter()
	if err != nil {
		log.Fatal(err)
		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
	})
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

func NewDb(lifecycle fx.Lifecycle, e *services.Env) (*sqlx.DB, error) {
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("closing DB connection")
			return db.Close()
		},
	})

	return db, nil
}

func NewRedisCli(lifecycle fx.Lifecycle, e *services.Env) (*redis.Client, error) {
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

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("closing redis connection")
			return redisCli.Close()
		},
	})
	return redisCli, nil
}
