package services

type CacheRepository interface {
	Get(key string) (value string, exists bool, error error)
	Set(key, value string, ttl int) error
}

type Cache struct {
	repository CacheRepository
}

func NewCacheService(r CacheRepository) *Cache {
	return &Cache{repository: r}
}

func (c *Cache) Get(key string) (value string, exists bool, error error) {
	return c.repository.Get(key)
}

func (c *Cache) Set(key, value string, ttl int) error {
	return c.repository.Set(key, value, ttl)
}
