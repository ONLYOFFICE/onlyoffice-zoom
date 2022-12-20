package rpc

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/messaging"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/wrapper"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/registry"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/trace"
	"github.com/go-micro/plugins/v4/wrapper/breaker/hystrix"
	rlimiter "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	"github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	uber "go.uber.org/ratelimit"
)

type Service struct {
	micro.Service
}

// NewService Initializes an http service with options.
func NewService(opts ...Option) Service {
	var tracer *oteltrace.TracerProvider
	var err error
	var wrappers []server.HandlerWrapper = make([]server.HandlerWrapper, 0, 2)

	options := newOptions(opts...)

	if options.Server == nil {
		log.Fatal("rpc service expects to have an initialized rpc server")
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

	if options.Limits > 0 {
		wrappers = append(wrappers, rlimiter.NewHandlerWrapper(int(options.Limits), uber.Per(1*time.Second)))
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
		wrappers = append(wrappers, wrapper.TracePropagationHandlerWrapper)
	}

	hystrix.ConfigureDefault(hystrix.CommandConfig{Timeout: 2500})
	service := micro.NewService(
		micro.Name(strings.Join([]string{options.Namespace, options.Name}, ":")),
		micro.Version(strconv.Itoa(options.Version)),
		micro.Context(options.Context),
		micro.Server(server.NewServer(
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
			opentelemetry.NewClientWrapper(opentelemetry.WithTraceProvider(otel.GetTracerProvider())),
		),
		micro.WrapSubscriber(opentelemetry.NewSubscriberWrapper(opentelemetry.WithTraceProvider(otel.GetTracerProvider()))),
		micro.WrapHandler(wrappers...),
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

	if len(options.Server.BuildMessageHandlers()) > 0 {
		for _, entry := range options.Server.BuildMessageHandlers() {
			if entry.Handler != nil && entry.Topic != "" {
				if entry.Queue != "" {
					if err := micro.RegisterSubscriber(
						entry.Topic, service.Server(), entry.Handler, server.SubscriberQueue(entry.Queue),
					); err != nil {
						log.Fatalf("could not register a new subscriber with topic %s. Reason: %s", entry.Topic, err.Error())
					}
				} else {
					if err := micro.RegisterSubscriber(entry.Topic, service.Server(), entry.Handler); err != nil {
						log.Fatalf("could not register a new subscriber with topic %s. Reason: %s", entry.Topic, err.Error())
					}
				}
			}
		}
	}

	for _, handler := range options.Server.BuildHandlers(service.Client()) {
		if err := micro.RegisterHandler(service.Server(), handler); err != nil {
			log.Fatalf("could not initialize rpc handlers: %s", err.Error())
		}
	}

	return Service{service}
}
