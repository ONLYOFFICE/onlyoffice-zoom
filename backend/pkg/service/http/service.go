package http

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/messaging"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/cors"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/registry"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/trace"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	mserver "github.com/go-micro/plugins/v4/server/http"
	"github.com/go-micro/plugins/v4/wrapper/breaker/hystrix"
	"github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
)

type Service struct {
	micro.Service
}

// NewService Initializes an http service with options.
func NewService(opts ...Option) Service {
	var tracer *oteltrace.TracerProvider
	var err error

	options := newOptions(opts...)

	if options.Server == nil {
		log.Fatal("http service expects to have an initialized http server")
	}

	registry := registry.NewRegistry(
		registry.WithAddresses(options.RegistryOptions.Addresses...),
		registry.WithSecure(options.RegistryOptions.Secure),
		registry.WithRegisryType(options.RegistryOptions.RegistryType),
		registry.WithCacheTTL(options.RegistryOptions.CacheTTL),
	)

	broker := messaging.NewBroker(
		registry,
		messaging.WithBrokerType(options.BrokerOptions.BrokerType),
		messaging.WithAddrs(options.BrokerOptions.Addrs...),
		messaging.WithSecure(options.BrokerOptions.Secure),
		messaging.WithContext(options.Context),
	)

	if err = broker.Init(); err != nil {
		log.Fatalf("could not initialize a new broker instance: %s", err.Error())
	}

	if err = broker.Connect(); err != nil {
		log.Fatalf("broker connection error: %s", err.Error())
	}

	if options.Tracer.Enable {
		if tracer, err = trace.NewTracer(
			trace.WithName(options.Tracer.Name),
			trace.WithAddress(options.Tracer.Address),
			trace.WithFractionRatio(options.Tracer.FractionRatio),
			trace.WithTracerType(options.Tracer.TracerType),
		); err != nil {
			log.Fatalf("could not initialize a new tracer instance: %s", err.Error())
		}
	}

	hystrix.ConfigureDefault(hystrix.CommandConfig{Timeout: 2500})
	service := micro.NewService(
		micro.Name(strings.Join([]string{options.Namespace, options.Name}, ":")),
		micro.Version(strconv.Itoa(options.Version)),
		micro.Context(options.Context),
		micro.Server(mserver.NewServer(
			server.Name(strings.Join([]string{options.Namespace, options.Name}, ":")),
			server.Address(options.Address),
		)),
		micro.Registry(registry),
		micro.Broker(broker),
		micro.Client(client.NewClient(
			client.Registry(registry),
			client.Broker(broker),
		)),
		micro.WrapClient(
			roundrobin.NewClientWrapper(),
			hystrix.NewClientWrapper(),
		),
		micro.WrapCall(opentelemetry.NewCallWrapper(opentelemetry.WithTraceProvider(otel.GetTracerProvider()))),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
		micro.AfterStop(func() error {
			if tracer != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := tracer.Shutdown(ctx); err != nil {
					return err
				}

				return nil
			}

			return nil
		}),
	)

	if options.IPLimits > 0 {
		options.Server.ApplyMiddleware(middleware.NewRateLimiter(options.IPLimits, 1*time.Second, middleware.WithKeyFuncIP))
	}

	if options.Limits > 0 {
		options.Server.ApplyMiddleware(middleware.NewRateLimiter(options.Limits, 1*time.Second, middleware.WithKeyFuncAll))
	}

	options.Server.ApplyMiddleware(
		middleware.Log(options.Logger),
		chimiddleware.RealIP,
		chimiddleware.RequestID,
		chimiddleware.StripSlashes,
		middleware.Secure,
		middleware.Version(strconv.Itoa(options.Version)),
		middleware.Cors(
			cors.WithAllowCredentials(options.CORS.AllowCredentials),
			cors.WithAllowedHeaders(options.CORS.AllowedHeaders...),
			cors.WithAllowedMethods(options.CORS.AllowedMethods...),
			cors.WithAllowedOrigins(options.CORS.AllowedOrigins...),
		),
	)

	if options.Tracer.Enable {
		options.Server.ApplyMiddleware(
			middleware.TracePropagationMiddleware,
			middleware.LogTraceMiddleware,
		)
	}

	if err := micro.RegisterHandler(
		service.Server(),
		options.Server.NewHandler(service.Options().Client),
	); err != nil {
		log.Fatal("could not register http handler: ", err)
	}

	return Service{service}
}
