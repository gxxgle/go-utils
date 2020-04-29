package sync

import (
	"time"

	"github.com/bsm/redislock"
)

type Options struct {
	ttl     time.Duration           // redislock ttl
	rlRetry redislock.RetryStrategy // redislock retry strategy
}

type Option func(*Options)

func TTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.ttl = ttl
	}
}

func Retry(retry redislock.RetryStrategy) Option {
	return func(o *Options) {
		o.rlRetry = retry
	}
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		ttl:     time.Minute,
		rlRetry: redislock.LinearBackoff(time.Millisecond * 100),
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}
