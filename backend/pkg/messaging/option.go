package messaging

import (
	"context"
)

type BrokerType int

var (
	RabbitMQ BrokerType = 1
	NATS     BrokerType = 2
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	BrokerType BrokerType
	Addrs      []string
	Secure     bool
	Context    context.Context
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		Secure:  false,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithBrokerType sets broker type
func WithBrokerType(val BrokerType) Option {
	return func(o *Options) {
		o.BrokerType = val
	}
}

// WithAddrs sets broker addresses.
func WithAddrs(val ...string) Option {
	return func(o *Options) {
		if len(val) > 0 {
			o.Addrs = val
		}
	}
}

// WithSecure sets broker secure flag
func WithSecure(val bool) Option {
	return func(o *Options) {
		o.Secure = val
	}
}

// WithContext sets broker context
func WithContext(val context.Context) Option {
	return func(o *Options) {
		if val != nil {
			o.Context = val
		}
	}
}
