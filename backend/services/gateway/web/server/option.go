package server

import (
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options
type Options struct {
	Logger        log.Logger
	ClientID      string
	ClientSecret  string
	WebhookSecret string
	RedirectURI   string
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

// WithLogger sets logger option.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

func WithClientID(val string) Option {
	return func(o *Options) {
		o.ClientID = val
	}
}

func WithClientSecret(val string) Option {
	return func(o *Options) {
		o.ClientSecret = val
	}
}

func WithWebhookSecret(val string) Option {
	return func(o *Options) {
		o.WebhookSecret = val
	}
}

func WithRedirectURI(val string) Option {
	return func(o *Options) {
		o.RedirectURI = val
	}
}
