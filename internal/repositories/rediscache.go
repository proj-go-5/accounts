package store

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/proj-go-5/accounts/internal/services"
)

type RedisCache struct {
	cli *redis.Client
}

func NewRedisCacheRepository(cli *redis.Client) services.CacheRepository {
	return &RedisCache{cli: cli}
}

func (r *RedisCache) Get(key string) (value string, exists bool, error error) {
	cmd := r.cli.Get(key)
	value = cmd.Val()
	if value != "" {
		exists = true
	}
	return cmd.Val(), exists, nil
}

func (r *RedisCache) Set(key, value string, ttl time.Duration) error {

	r.cli.Set(key, value, ttl)
	return nil
}
