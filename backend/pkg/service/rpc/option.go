package rpc

import (
	"context"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/messaging"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/registry"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/trace"
	"go-micro.dev/v4/client"
)

type RPCMessageHandler struct {
	Topic   string
	Queue   string
	Handler interface{}
}

type RPCEngine interface {
	BuildMessageHandlers() []RPCMessageHandler
	BuildHandlers(client client.Client) []interface{}
}

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	Namespace                           string
	Name                                string
	Version                             int
	Address                             string
	Limits                              uint64
	CircuitBreakerTimeout               int
	CircuitBreakerMaxConcurrent         int
	CircuitBreakerVolumeThreshold       int
	CircuitBreakerSleepWindow           int
	CircuitBreakerErrorPercentThreshold int
	Server                              RPCEngine
	Logger                              log.Logger
	Tracer                              trace.Options
	BrokerOptions                       messaging.Options
	RegistryOptions                     registry.Options
	Context                             context.Context
}

// newOptions initializes the options.
func newOptions(opts ...Option) Options {
	opt := Options{
		Namespace: "anonymous-namespace.go-micro",
		Name:      "anonymous-rpc",
		Version:   0,
		Address:   ":7070",
		Logger:    log.NewDefaultLogger(),
		RegistryOptions: registry.Options{
			CacheTTL:     10 * time.Second,
			RegistryType: 4,
		},
		Context: context.Background(),
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
		if val >= 0 {
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

// WithLimits sets rate limiter limits.
func WithLimits(val uint64) Option {
	return func(o *Options) {
		if val > 0 {
			o.Limits = val
		}
	}
}

// WithCircuitBreakerTimeout sets hystrix timeout
func WithCircuitBreakerTimeout(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.CircuitBreakerTimeout = val
		}
	}
}

// WithCircuitBreakerMaxConcurrent sets hystrix max concurrency level
func WithCircuitBreakerMaxConcurrent(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.CircuitBreakerMaxConcurrent = val
		}
	}
}

// WithCircuitBreakerVolumeThreshold sets hystrix threshold
func WithCircuitBreakerVolumeThreshold(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.CircuitBreakerVolumeThreshold = val
		}
	}
}

// WithCircuitBreakerSleepWindow sets hystrix sleep window
func WithCircuitBreakerSleepWindow(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.CircuitBreakerSleepWindow = val
		}
	}
}

// WithCircuitBreakerErrorPercentThreshold sets hystrix error threshold
func WithCircuitBreakerErrorPercentThreshold(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.CircuitBreakerErrorPercentThreshold = val
		}
	}
}

// WithBrokerOptions sets messaging configuration
func WithBrokerOptions(val messaging.Options) Option {
	return func(o *Options) {
		o.BrokerOptions = val
	}
}

// WithRPCServer sets RPCServer config
func WithRPC(val RPCEngine) Option {
	return func(o *Options) {
		if val != nil {
			o.Server = val
		}
	}
}

// WithTracer turns on/off distributed tracing
func WithTracer(val trace.Options) Option {
	return func(o *Options) {
		o.Tracer = val
	}
}

// WithLogger sets server logger.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		if val != nil {
			o.Logger = val
		}
	}
}

// WithRegistryOptions passes registry options to the service initialization.
func WithRegistryOptions(val registry.Options) Option {
	return func(o *Options) {
		if len(val.Addresses) > 0 {
			o.RegistryOptions.Addresses = val.Addresses
		}

		if val.CacheTTL > 0 {
			o.RegistryOptions.CacheTTL = val.CacheTTL
		}

		if val.RegistryType > 0 {
			o.RegistryOptions.RegistryType = val.RegistryType
		}
	}
}

// WithContext sets the context option.
func WithContext(val context.Context) Option {
	return func(o *Options) {
		if val != nil {
			o.Context = val
		}
	}
}
