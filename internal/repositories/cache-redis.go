package store

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	cli *redis.Client
}

func NewRedisCacheRepository(r *redis.Client) *RedisCache {
	return &RedisCache{cli: r}
}

func (r *RedisCache) Get(key string) (value string, exists bool, error error) {
	cmd := r.cli.Get(key)
	value = cmd.Val()
	if value != "" {
		exists = true
	}
	return cmd.Val(), exists, nil
}

func (r *RedisCache) Set(key, value string, ttl int) error {

	r.cli.Set(key, value, time.Duration(ttl))
	return nil
}
