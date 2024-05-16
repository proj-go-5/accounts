package store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/proj-go-5/accounts/internal/services"
)

type RedisCache struct {
	Cli *redis.Client
}

func NewRedisCacheRepository(e *services.Env) (*RedisCache, error) {
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

	return &RedisCache{Cli: redisCli}, nil
}

func (r *RedisCache) Get(key string) (value string, exists bool, error error) {
	cmd := r.Cli.Get(key)
	value = cmd.Val()
	if value != "" {
		exists = true
	}
	return cmd.Val(), exists, nil
}

func (r *RedisCache) Set(key, value string, ttl time.Duration) error {

	r.Cli.Set(key, value, ttl)
	return nil
}
