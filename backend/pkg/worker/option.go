package worker

import (
	"crypto/tls"
	"strings"
	"time"
)

type Option func(*Options)

// Options defines the available options.
type Options struct {
	MaxActive         int
	MaxIdle           int
	MaxConcurrency    uint
	RedisNamespace    string
	RedisAddress      string
	RedisUsername     string
	RedisPassword     string
	RedisDatabase     int
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
	TLSConfig         *tls.Config
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		MaxActive:         5,
		MaxIdle:           5,
		MaxConcurrency:    4,
		RedisNamespace:    "common",
		RedisAddress:      "0.0.0.0:6379",
		RedisDatabase:     0,
		RedisReadTimeout:  2 * time.Second,
		RedisWriteTimeout: 3 * time.Second,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithMaxActive(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.MaxActive = val
		}
	}
}

func WithMaxIdle(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.MaxIdle = val
		}
	}
}

func WithMaxConcurrency(val uint) Option {
	return func(o *Options) {
		if val > 0 {
			o.MaxConcurrency = val
		}
	}
}

func WithRedisNamespace(val string) Option {
	return func(o *Options) {
		o.RedisNamespace = val
	}
}

func WithRedisAddress(val string) Option {
	return func(o *Options) {
		o.RedisAddress = val
	}
}

func WithRedisUsername(val string) Option {
	return func(o *Options) {
		o.RedisUsername = val
	}
}

func WithRedisPassword(val string) Option {
	return func(o *Options) {
		o.RedisPassword = strings.TrimSpace(val)
	}
}

func WithRedisDatabase(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.RedisDatabase = val
		}
	}
}

func WithTLSConfig(val *tls.Config) Option {
	return func(o *Options) {
		if val != nil {
			o.TLSConfig = val
		}
	}
}
