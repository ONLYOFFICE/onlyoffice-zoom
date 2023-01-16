package registry

import "time"

type RegistryType int

var (
	Kubernetes RegistryType = 1
	Consul     RegistryType = 2
	Etcd       RegistryType = 3
	MDNS       RegistryType = 4
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	Addresses    []string
	CacheTTL     time.Duration
	RegistryType RegistryType
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		CacheTTL: 10 * time.Second,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithAddresses sets registry addresses.
func WithAddresses(val ...string) Option {
	return func(o *Options) {
		if len(val) > 0 {
			o.Addresses = val
		}
	}
}

// WithCacheTTL sets registry lookup cache.
func WithCacheTTL(val time.Duration) Option {
	return func(o *Options) {
		if val > 0 {
			o.CacheTTL = val
		}
	}
}

// WithRegistryType sets registry type.
func WithRegisryType(val RegistryType) Option {
	return func(o *Options) {
		if val > 0 {
			o.RegistryType = val
		}
	}
}
