package store

import "sync"

type MemoryCache struct {
	mx    sync.Mutex
	store map[string]string
}

func NewMemoryCacheRepository() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]string),
	}
}

func (c *MemoryCache) Get(key string) (value string, exists bool, error error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	value, exists = c.store[key]
	return value, exists, nil
}

func (c *MemoryCache) Set(key, value string, ttl int) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.store[key] = value
	return nil
}
