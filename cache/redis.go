package cache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
)

// RedisConfig is config struct of redis.
type RedisConfig struct {
	URL      string `json:"url"`
	Password string `json:"password"`
	Retries  int    `json:"retries"`
}

type RedisCacher struct {
	*redis.Client
}

func InitRedis(cfg *RedisConfig) error {
	c, err := NewRedisCacher(cfg)
	if err != nil {
		return err
	}

	DefaultCacher = c
	return nil
}

func NewRedisCacher(cfg *RedisConfig) (Cacher, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:        cfg.URL,
		Password:    cfg.Password,
		DB:          0,
		DialTimeout: time.Second * 5,
	})

	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}

	opt := cli.Options()
	opt.MaxRetries = cfg.Retries
	return &RedisCacher{cli}, nil
}

// Set value marshal to json
func (c *RedisCacher) Set(key string, value interface{}, expiration time.Duration) error {
	str, err := json.MarshalToString(value)
	if err != nil {
		return err
	}

	return c.Client.Set(key, str, expiration).Err()
}

// Get value unmarshal from json
func (c *RedisCacher) Get(key string, value interface{}) error {
	bs, err := c.Client.Get(key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, value)
}

// HSet value to json
func (c *RedisCacher) HSet(key, field string, value interface{}) error {
	str, err := json.MarshalToString(value)
	if err != nil {
		return err
	}

	return c.Client.HSet(key, field, str).Err()
}

// HGet value from json
func (c *RedisCacher) HGet(key, field string, value interface{}) error {
	bs, err := c.Client.HGet(key, field).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, value)
}

func (c *RedisCacher) Close() {
	if err := c.Client.Close(); err != nil {
		log.Errorw("redis close error", "err", err)
	}
}
