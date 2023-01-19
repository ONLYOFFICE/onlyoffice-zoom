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
	Logger     log.Logger
	Config     config.HttpServer
	Callback   shared.CallbackConfig
	Onlyoffice shared.OnlyofficeConfig
	Tracer     config.Tracer
	Broker     config.Broker
	Registry   config.Registry
	Worker     config.WorkerConfig
	Context    context.Context
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

// WithConfig sets http server config.
func WithConfig(val config.HttpServer) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// WithOnlyofficeConfig sets doc server's info
func WithOnlyoffice(val shared.OnlyofficeConfig) Option {
	return func(o *Options) {
		o.Onlyoffice = val
	}
}

// WithCallback sets callback handler's config.
func WithCallback(val shared.CallbackConfig) Option {
	return func(o *Options) {
		o.Callback = val
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

func WithWorker(val config.WorkerConfig) Option {
	return func(o *Options) {
		o.Worker = val
	}
}
