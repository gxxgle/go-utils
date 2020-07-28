package sync

import (
	"github.com/bsm/redislock"
	"github.com/gxxgle/go-utils/cache"
	"github.com/gxxgle/go-utils/log"
)

type redisMutexer struct {
	options *Options
	cacher  *cache.RedisCacher
	locker  *redislock.Client
}

type redisMutex struct {
	key     string
	mutexer *redisMutexer
	lock    *redislock.Lock
}

func InitRedis(cfg *cache.RedisConfig, opts ...Option) error {
	mu, err := NewRedisMutexer(cfg, opts...)
	if err != nil {
		return err
	}

	DefaultMutexer = mu
	return nil
}

func NewRedisMutexer(cfg *cache.RedisConfig, opts ...Option) (Mutexer, error) {
	cfg.Retries = 0
	cacher, err := cache.NewRedisCacher(cfg)
	if err != nil {
		return nil, err
	}

	out := &redisMutexer{
		options: newOptions(opts...),
		cacher:  cacher.(*cache.RedisCacher),
		locker:  redislock.New(cacher.(*cache.RedisCacher)),
	}

	return out, nil
}

func (m *redisMutexer) NewMutex(key string) Mutex {
	return &redisMutex{
		key:     key,
		mutexer: m,
	}
}

func (m *redisMutexer) Close() {
	m.cacher.Close()
}

func (m *redisMutex) Lock() {
	var (
		err error
		opt = m.mutexer.options
	)

	for {
		m.lock, err = m.mutexer.locker.Obtain(m.key, opt.ttl, &redislock.Options{RetryStrategy: opt.rlRetry})
		if err == nil {
			return
		}

		log.L.WithError(err).WithField("key", m.key).Error("go-utils redislock lock error")
	}
}

func (m *redisMutex) Unlock() {
	if err := m.lock.Release(); err != nil {
		log.L.WithError(err).WithField("key", m.key).Error("go-utils redislock unlock error")
	}
}
