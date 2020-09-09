package cache

import (
	"time"

	"github.com/gxxgle/go-utils/json"

	"github.com/go-redis/redis/v7"
	"github.com/phuslu/log"
)

// RedisConfig is config struct of redis.
type RedisConfig struct {
	Network    string `json:"network" yaml:"network" validate:"default=tcp,oneof=tcp unix"`
	Addr       string `json:"addr" yaml:"addr" validate:"required"`
	Password   string `json:"password" yaml:"password"`
	DB         int    `json:"db" yaml:"db"`
	MaxRetries int    `json:"max_retries" yaml:"max_retries"`
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
		Network:    cfg.Network,
		Addr:       cfg.Addr,
		Password:   cfg.Password,
		DB:         cfg.DB,
		MaxRetries: cfg.MaxRetries,
	})

	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}

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

// Delete cache by key
func (c *RedisCacher) Delete(key ...string) error {
	return c.Client.Del(key...).Err()
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
		log.Error().Err(err).Msg("redis close")
	}
}
