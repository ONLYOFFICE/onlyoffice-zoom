package service

import (
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/messaging"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/cors"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/registry"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/http"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/trace"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/worker"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/server"
)

// NewService used to register a new go-micro service
func NewService(opts ...Option) (*http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.WithNamespace(options.Config.Namespace),
		http.WithName(options.Config.Name),
		http.WithVersion(options.Config.Version),
		http.WithAddress(options.Config.Address),
		http.WithLimits(options.Config.Resilience.RateLimiter.Limit),
		http.WithIPLimits(options.Config.Resilience.RateLimiter.IPLimit),
		http.WithCircuitBreakerVolumeThreshold(options.Config.Resilience.CircuitBreaker.VolumeThreshold),
		http.WithCircuitBreakerTimeout(options.Config.Resilience.CircuitBreaker.Timeout),
		http.WithCircuitBreakerSleepWindow(options.Config.Resilience.CircuitBreaker.SleepWindow),
		http.WithCircuitBreakerMaxConcurrent(options.Config.Resilience.CircuitBreaker.MaxConcurrent),
		http.WithCircuitBreakerErrorPercentThreshold(options.Config.Resilience.CircuitBreaker.ErrorPercentThreshold),
		http.WithServer(
			server.NewServer(
				server.WithNamespace(options.Config.Namespace),
				server.WithLogger(options.Logger),
				server.WithDocSecret(options.Onlyoffice.DocSecret),
				server.WithMaxSize(options.Callback.MaxSize),
				server.WithUploadTimeout(options.Callback.UploadTimeout),
				server.WithWorker(worker.NewOptions(
					worker.WithMaxActive(options.Worker.MaxActive),
					worker.WithMaxIdle(options.Worker.MaxIdle),
					worker.WithMaxConcurrency(options.Worker.MaxConcurrency),
					worker.WithRedisNamespace(options.Worker.RedisNamespace),
					worker.WithRedisAddress(options.Worker.RedisAddress),
					worker.WithRedisUsername(options.Worker.RedisUsername),
					worker.WithRedisPassword(options.Worker.RedisPassword),
					worker.WithRedisDatabase(options.Worker.RedisDatabase),
				)),
			),
		),
		http.WithLogger(options.Logger),
		http.WithTracer(trace.NewOptions(
			trace.WithEnable(options.Tracer.Enable),
			trace.WithName(strings.Join([]string{options.Config.Namespace, options.Config.Name}, ":")),
			trace.WithAddress(options.Tracer.Address),
			trace.WithFractionRatio(options.Tracer.FractionRatio),
			trace.WithTracerType(trace.TracerType(options.Tracer.TracerType)),
		)),
		http.WithCORS(cors.Options{
			AllowedOrigins:   options.Config.CORS.AllowedOrigins,
			AllowedMethods:   options.Config.CORS.AllowedMethods,
			AllowedHeaders:   options.Config.CORS.AllowedHeaders,
			AllowCredentials: options.Config.CORS.AllowedCredentials,
		}),
		http.WithBrokerOptions(messaging.NewOptions(
			messaging.WithAddrs(options.Broker.Addrs...),
			messaging.WithBrokerType(messaging.BrokerType(options.Broker.Type)),
			messaging.WithDisableAutoAck(options.Broker.DisableAutoAck),
			messaging.WithDurable(options.Broker.Durable),
			messaging.WithRequeueOnError(options.Broker.RequeueOnError),
			messaging.WithAckOnSuccess(options.Broker.AckOnSuccess),
		)),
		http.WithRegistryOptions(registry.Options{
			Addresses:    options.Registry.Addresses,
			CacheTTL:     options.Registry.CacheTTL,
			RegistryType: registry.RegistryType(options.Registry.RegistryType),
		}),
		http.WithContext(options.Context),
	)

	return &service, nil
}
