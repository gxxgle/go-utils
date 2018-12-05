package sync

import (
	"time"

	"github.com/bsm/redis-lock"
	"github.com/gxxgle/go-utils/cache"
	"github.com/gxxgle/go-utils/log"
)

type redisMutexer struct {
	cacher *cache.RedisCacher
	option *lock.Options
}

type redisMutex struct {
	*lock.Locker
	key string
}

func InitRedis(cfg *cache.RedisConfig, opt ...*lock.Options) error {
	mu, err := NewRedisMutexer(cfg, opt...)
	if err != nil {
		return err
	}

	DefaultMutexer = mu
	return nil
}

func NewRedisMutexer(cfg *cache.RedisConfig, opt ...*lock.Options) (Mutexer, error) {
	cfg.Retries = 0
	cacher, err := cache.NewRedisCacher(cfg)
	if err != nil {
		return nil, err
	}

	out := &redisMutexer{
		cacher: cacher.(*cache.RedisCacher),
		option: &lock.Options{
			LockTimeout: time.Second * 60,
			RetryDelay:  time.Millisecond * 20,
		},
	}

	if len(opt) > 0 {
		out.option = opt[0]
	}

	out.option.RetryCount = int(out.option.LockTimeout/out.option.RetryDelay) + 1
	return out, nil
}

func (m *redisMutexer) NewMutex(key string) Mutex {
	return &redisMutex{
		Locker: lock.New(m.cacher, key, m.option),
		key:    key,
	}
}

func (m *redisMutexer) Close() {
	m.cacher.Close()
}

func (m *redisMutex) Lock() {
	for {
		_, err := m.Locker.Lock()
		if err == nil {
			return
		}

		log.Errorw("redis sync lock error", "key", m.key, "err", err)
	}
}

func (m *redisMutex) Unlock() {
	if err := m.Locker.Unlock(); err != nil {
		log.Errorw("redis sync unlock error", "key", m.key, "err", err)
	}
}
