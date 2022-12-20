package adapter

import (
	"crypto/tls"
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
)

type Option func(*Options)

// Options defines the available options.
type Options struct {
	RedisAddresses []string
	RedisUsername  string
	RedisPassword  string
	BufferSize     int
	TLSConfig      *tls.Config
	Logger         log.Logger
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		BufferSize:     100,
		Logger:         log.NewDefaultLogger(),
		RedisAddresses: []string{"0.0.0.0:6379"},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithRedisAddresses(val []string) Option {
	return func(o *Options) {
		o.RedisAddresses = val
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

func WithBufferSize(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.BufferSize = val
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

func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}
