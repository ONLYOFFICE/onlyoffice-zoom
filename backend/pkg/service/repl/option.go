package repl

import (
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/cors"
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	Namespace string
	Name      string
	Version   int
	Address   string
	Debug     bool
	CORS      cors.Options
}

// newOptions initializes the options.
func newOptions(opts ...Option) Options {
	opt := Options{
		Namespace: "anonymous-namespace.go-micro",
		Name:      "anonymous-repl",
		Version:   0,
		Address:   ":8082",
		Debug:     true,
		CORS:      cors.NewOptions(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithNamespace sets the namespace option.
func WithNamespace(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Namespace = val
		}
	}
}

// WithName sets the name option.
func WithName(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Name = val
		}
	}
}

// WithVersion sets the version option.
func WithVersion(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.Version = val
		}
	}
}

// WithAddress sets the address option.
func WithAddress(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Address = val
		}
	}
}

// WithDebug sets pprof flag.
func WithDebug(val bool) Option {
	return func(o *Options) {
		o.Debug = val
	}
}
