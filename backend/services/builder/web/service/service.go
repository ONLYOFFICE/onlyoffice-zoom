package service

import (
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/messaging"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/registry"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/rpc"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/trace"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server"
)

// NewService used to register a new go-micro service
func NewService(opts ...Option) (*rpc.Service, error) {
	options := newOptions(opts...)

	service := rpc.NewService(
		rpc.WithNamespace(options.Config.Namespace),
		rpc.WithName(options.Config.Name),
		rpc.WithVersion(options.Config.Version),
		rpc.WithAddress(options.Config.Address),
		rpc.WithLimits(options.Config.Resilience.RateLimiter.Limit),
		rpc.WithCircuitBreakerTimeout(options.Config.Resilience.CircuitBreaker.Timeout),
		rpc.WithCircuitBreakerVolumeThreshold(options.Config.Resilience.CircuitBreaker.VolumeThreshold),
		rpc.WithCircuitBreakerSleepWindow(options.Config.Resilience.CircuitBreaker.SleepWindow),
		rpc.WithCircuitBreakerMaxConcurrent(options.Config.Resilience.CircuitBreaker.MaxConcurrent),
		rpc.WithCircuitBreakerErrorPercentThreshold(options.Config.Resilience.CircuitBreaker.ErrorPercentThreshold),
		rpc.WithRPC(server.NewConfigRPCServer(
			server.WithClientID(options.Zoom.ClientID),
			server.WithClientSecret(options.Zoom.ClientSecret),
			server.WithDocSecret(options.Onlyoffice.DocSecret),
			server.WithCallbackURL(options.Onlyoffice.CallbackURL),
			server.WithLogger(options.Logger),
			server.WithRedis(options.Redis),
		)),
		rpc.WithLogger(options.Logger),
		rpc.WithTracer(trace.NewOptions(
			trace.WithEnable(options.Tracer.Enable),
			trace.WithName(strings.Join([]string{options.Config.Namespace, options.Config.Name}, ":")),
			trace.WithAddress(options.Tracer.Address),
			trace.WithFractionRatio(options.Tracer.FractionRatio),
			trace.WithTracerType(trace.TracerType(options.Tracer.TracerType)),
		)),
		rpc.WithBrokerOptions(messaging.NewOptions(
			messaging.WithAddrs(options.Broker.Addrs...),
			messaging.WithBrokerType(messaging.BrokerType(options.Broker.Type)),
			messaging.WithContext(options.Context),
		)),
		rpc.WithRegistryOptions(registry.Options{
			Addresses:    options.Registry.Addresses,
			CacheTTL:     options.Registry.CacheTTL,
			RegistryType: registry.RegistryType(options.Registry.RegistryType),
		}),
		rpc.WithContext(options.Context),
	)

	return &service, nil
}
