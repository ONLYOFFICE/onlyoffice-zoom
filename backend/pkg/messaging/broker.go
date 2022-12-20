package messaging

import (
	"github.com/go-micro/plugins/v4/broker/memory"
	"github.com/go-micro/plugins/v4/broker/nats"
	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/registry"
)

// NewBroker create a broker instance based on BrokerType value
func NewBroker(registry registry.Registry, opts ...Option) broker.Broker {
	options := NewOptions(opts...)

	bo := []broker.Option{
		broker.Addrs(options.Addrs...),
		broker.Registry(registry),
		broker.Secure(options.Secure),
	}

	var b broker.Broker
	switch options.BrokerType {
	case RabbitMQ:
		b = rabbitmq.NewBroker(bo...)
	case NATS:
		b = nats.NewBroker(bo...)
	default:
		b = memory.NewBroker(bo...)
	}

	return b
}
