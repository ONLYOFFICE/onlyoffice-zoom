package server

import (
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared"
)

type Option func(*Options)

// Options defines the available options.
type Options struct {
	ClientID     string
	ClientSecret string
	DocSecret    string
	CallbackURL  string
	Logger       log.Logger
	Redis        shared.RedisConfig
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		Logger:    log.NewDefaultLogger(),
		DocSecret: "secret",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithClientID(val string) Option {
	return func(o *Options) {
		o.ClientID = strings.TrimSpace(val)
	}
}

func WithClientSecret(val string) Option {
	return func(o *Options) {
		o.ClientSecret = strings.TrimSpace(val)
	}
}

func WithDocSecret(val string) Option {
	return func(o *Options) {
		o.DocSecret = strings.TrimSpace(val)
	}
}

func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

func WithRedis(val shared.RedisConfig) Option {
	return func(o *Options) {
		o.Redis = val
	}
}

func WithCallbackURL(val string) Option {
	return func(o *Options) {
		o.CallbackURL = val
	}
}
