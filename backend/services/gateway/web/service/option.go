package service

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options.
type Options struct {
	Config   config.HttpServer
	Tracer   config.Tracer
	Broker   config.Broker
	Registry config.Registry
	Zoom     shared.ZoomConfig
	Cache    config.Cache
	Logger   log.Logger
	Context  context.Context
}

// newOptions initializes the options.
func newOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
		Logger:  log.NewDefaultLogger(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithConfig sets http server config.
func WithConfig(val config.HttpServer) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// WithTracer sets http tracer config
func WithTracer(val config.Tracer) Option {
	return func(o *Options) {
		o.Tracer = val
	}
}

// WithBroker sets http broker config
func WithBroker(val config.Broker) Option {
	return func(o *Options) {
		o.Broker = val
	}
}

// WithRegistry sets http server registry.
func WithRegistry(val config.Registry) Option {
	return func(o *Options) {
		o.Registry = val
	}
}

func WithZoomConfig(val shared.ZoomConfig) Option {
	return func(o *Options) {
		o.Zoom = val
	}
}

// WithLogger sets http server logger.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

// WithContext sets http server context.
func WithContext(val context.Context) Option {
	return func(o *Options) {
		if val != nil {
			o.Context = val
		}
	}
}

func WithCache(val config.Cache) Option {
	return func(o *Options) {
		o.Cache = val
	}
}
