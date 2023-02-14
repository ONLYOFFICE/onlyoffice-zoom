package server

import (
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/worker"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options
type Options struct {
	Namespace     string
	Logger        log.Logger
	Worker        worker.WorkerOptions
	DocSecret     string
	MaxSize       int64
	UploadTimeout int
}

// newOptions initializes the options.
func newOptions(opts ...Option) Options {
	options := Options{
		Namespace:     "onlyoffice",
		Logger:        log.NewDefaultLogger(),
		DocSecret:     "secret",
		MaxSize:       2100000,
		UploadTimeout: 10,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithNamespace(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Namespace = val
		}
	}
}

// WithLogger sets logger option.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

func WithDocSecret(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.DocSecret = val
		}
	}
}

func WithMaxSize(val int64) Option {
	return func(o *Options) {
		o.MaxSize = val
	}
}

func WithWorker(val worker.WorkerOptions) Option {
	return func(o *Options) {
		o.Worker = val
	}
}

func WithUploadTimeout(val int) Option {
	return func(o *Options) {
		o.UploadTimeout = val
	}
}
