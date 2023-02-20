package server

import (
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options
type Options struct {
	ClientID     string
	ClientSecret string
	Persistence  string
	Logger       log.Logger
}

// newOptions initializes the options.
func newOptions(opts ...Option) Options {
	options := Options{
		Logger: log.NewDefaultLogger(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithClientID(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.ClientID = val
		}
	}
}

func WithClientSecret(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.ClientSecret = val
		}
	}
}

func WithPersistence(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Persistence = val
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
