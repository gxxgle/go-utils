package cache

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gxxgle/go-utils/json"
	"github.com/patrickmn/go-cache"
)

type MemoryCacher struct {
	Client *cache.Cache
}

func InitMemory() {
	DefaultCacher = NewMemoryCacher()
}

func NewMemoryCacher() Cacher {
	return &MemoryCacher{
		Client: cache.New(cache.NoExpiration, time.Minute*10),
	}
}

// Set value marshal to json
func (c *MemoryCacher) Set(key string, value interface{}, expiration time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.Client.Set(key, bs, expiration)
	return nil
}

// Get value unmarshal from json
func (c *MemoryCacher) Get(key string, value interface{}) error {
	v, ok := c.Client.Get(key)
	if !ok {
		return redis.Nil
	}

	bs, ok := v.([]byte)
	if !ok {
		c.Client.Delete(key)
		return redis.Nil
	}

	return json.Unmarshal(bs, value)
}

// Delete cache by key
func (c *MemoryCacher) Delete(key ...string) error {
	for _, k := range key {
		c.Client.Delete(k)
	}

	return nil
}

func (c *MemoryCacher) Close() {
}
