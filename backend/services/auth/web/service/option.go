package service

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	aconfig "github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options.
type Options struct {
	Logger      log.Logger
	Config      config.RPCServer
	Persistence aconfig.Persistence
	Tracer      config.Tracer
	Broker      config.Broker
	Registry    config.Registry
	Zoom        shared.ZoomConfig
	Cache       config.Cache
	Context     context.Context
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

// WithLogger sets rpc server logger.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

// WithContext sets rpc server context.
func WithContext(val context.Context) Option {
	return func(o *Options) {
		if val != nil {
			o.Context = val
		}
	}
}

// WithConfig sets rpc server config.
func WithConfig(val config.RPCServer) Option {
	return func(o *Options) {
		o.Config = val
	}
}

func WithPersistence(val aconfig.Persistence) Option {
	return func(o *Options) {
		o.Persistence = val
	}
}

// WithZoomConfig sets zoom configuration
func WithZoomConfig(val shared.ZoomConfig) Option {
	return func(o *Options) {
		o.Zoom = val
	}
}

// WithTracer sets rpc tracer config
func WithTracer(val config.Tracer) Option {
	return func(o *Options) {
		o.Tracer = val
	}
}

// WithBroker sets rpc broker config
func WithBroker(val config.Broker) Option {
	return func(o *Options) {
		o.Broker = val
	}
}

// WithRegistry sets rpc server registry.
func WithRegistry(val config.Registry) Option {
	return func(o *Options) {
		o.Registry = val
	}
}

func WithCache(val config.Cache) Option {
	return func(o *Options) {
		o.Cache = val
	}
}
