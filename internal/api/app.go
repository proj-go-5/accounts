package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/authorization"
	"go.uber.org/fx"
)

type App struct {
	providers  []interface{}
	envService *services.Env
}

func NewApp() (*App, error) {

	providers := []interface{}{
		newEnvService,
		services.NewAdminService,
		services.NewCacheService,
		services.NewHashService,
		newJwtService,
		services.NewAuthService,
		services.NewAppService,
		NewApi,
	}
	envService, err := newEnvService()
	if err != nil {
		return nil, err
	}
	return &App{providers: providers, envService: envService}, nil
}

func (a *App) Run() {
	err := a.addAdminRepositoryProvider()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = a.addCacheRepositoryProvider()
	if err != nil {
		log.Fatal(err)
		return
	}

	fx.New(
		fx.Provide(a.providers...),

		fx.Invoke(a.runServer),
	).Run()
}

func (a *App) runServer(lifecycle fx.Lifecycle, api *API, e *services.Env) {
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

func newEnvService() (*services.Env, error) {
	envService, err := services.NewEnvService(".env")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return envService, nil
}

func newJwtService(e *services.Env) (*authorization.JwtService, error) {
	jwtSecret := e.Get("JWT_SECRET", "secret")
	jwtExpiration, err := strconv.Atoi(e.Get("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, err
	}
	return authorization.NewJwtService(jwtSecret, jwtExpiration), nil
}

func newDb(lifecycle fx.Lifecycle, e *services.Env) (*sqlx.DB, error) {
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

func newRedisCli(lifecycle fx.Lifecycle, e *services.Env) (*redis.Client, error) {
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

func (a *App) addAdminRepositoryProvider() error {
	dbRepositoryTypeEnvVarName := "ACCOUNTS_DB_REPOSITORY_TYPE"

	dbRepositoryType := a.envService.Get(dbRepositoryTypeEnvVarName, "")
	if dbRepositoryType == "" {
		return fmt.Errorf("%v env variable should be provided", dbRepositoryTypeEnvVarName)
	}

	if dbRepositoryType != "postgres" && dbRepositoryType != "in_memory" {
		return fmt.Errorf("%v env varialbe could be with 'postgres' or 'in_memory' values only", dbRepositoryTypeEnvVarName)
	}

	if dbRepositoryType == "postgres" {
		a.providers = append(a.providers, newDb, store.NewAdminDBRepository)
	} else {
		a.providers = append(a.providers, store.NewAdminMemoryRepository)
	}

	log.Printf("Using %v=%v", dbRepositoryTypeEnvVarName, dbRepositoryType)
	return nil
}

func (a *App) addCacheRepositoryProvider() error {
	cacheRepositoryTypeEnvVarName := "ACCOUNTS_CACHE_REPOSITORY_TYPE"

	cacheRepositoryType := a.envService.Get(cacheRepositoryTypeEnvVarName, "")
	if cacheRepositoryType == "" {
		return fmt.Errorf("%v env variable should be provided", cacheRepositoryTypeEnvVarName)
	}

	if cacheRepositoryType != "redis" && cacheRepositoryType != "in_memory" {
		return fmt.Errorf("%v env varialbe could be with 'redis' or 'in_memory' values only", cacheRepositoryTypeEnvVarName)
	}

	if cacheRepositoryType == "redis" {
		a.providers = append(a.providers, newRedisCli, store.NewRedisCacheRepository)
	} else {
		a.providers = append(a.providers, store.NewMemoryCacheRepository)
	}
	log.Printf("Using %v=%v", cacheRepositoryTypeEnvVarName, cacheRepositoryType)

	return nil
}
