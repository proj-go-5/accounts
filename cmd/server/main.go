package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/proj-go-5/accounts/internal/api"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/env"
)

func main() {
	envService, err := env.NewEnvService(".env")
	if err != nil {
		log.Println(err)
		return
	}

	dbDataSource := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v",
		envService.Get("ACCOUNTS_DB_USER", "accouunts"),
		envService.Get("ACCOUNTS_DB_PASSWORD", "accouunts"),
		envService.Get("ACCOUNTS_DB_NAME", "accouunts"),
		envService.Get("ACCOUNTS_DB_HOST", "localhost"),
		envService.Get("ACCOUNTS_DB_PORT", "5432"),
		envService.Get("ACCOUNTS_DB_SSL_MODE", "disable"),
	)

	db, err := sqlx.Open("postgres", dbDataSource)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return
	}

	redisAddres := fmt.Sprintf("%v:%v",
		envService.Get("ACCOUNTS_REDIS_HOST", "localhost"),
		envService.Get("ACCOUNTS_REDIS_PORT", "6379"),
	)

	redisDb, err := strconv.Atoi(envService.Get("ACCOUNTS_REDIS_DB", "0"))
	if err != nil {
		log.Fatal(err)
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr:     redisAddres,
		Password: envService.Get("ACCOUNTS_REDIS_PASSWORD", ""),
		DB:       redisDb,
	})
	defer redisCli.Close()

	var cacheStore services.CacheRepository

	err = redisCli.Ping().Err()

	if err != nil {
		log.Fatal(err)
		cacheStore = store.NewMemoryCacheRepository()
	} else {
		cacheStore = store.NewRedisCacheRepository(redisCli)
	}

	userService := services.NewUserService(store.NewUserDBRepository(db))
	cacheService := services.NewCacheService(cacheStore)
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

	serverPort := envService.Get("ACCOUNTS_SERVER_PORT", "8080")

	fmt.Printf("Runing servier on %v port\n", serverPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
		log.Printf("Server run error: %s", err)
	}
}
