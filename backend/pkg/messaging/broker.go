package messaging

import (
	"github.com/go-micro/plugins/v4/broker/memory"
	"github.com/go-micro/plugins/v4/broker/nats"
	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/registry"
)

// NewBroker create a broker instance based on BrokerType value
func NewBroker(registry registry.Registry, opts ...Option) (broker.Broker, broker.SubscribeOptions) {
	// TODO: Additional brokers and options
	options := NewOptions(opts...)

	bo := []broker.Option{
		broker.Addrs(options.Addrs...),
		broker.Registry(registry),
	}

	var b broker.Broker
	var subOpts broker.SubscribeOptions

	switch options.BrokerType {
	case RabbitMQ:
		b = rabbitmq.NewBroker(bo...)

		opts := []broker.SubscribeOption{}
		if options.DisableAutoAck {
			opts = append(opts, broker.DisableAutoAck())
		}

		if options.AckOnSuccess {
			opts = append(opts, rabbitmq.AckOnSuccess())
		}

		if options.Durable {
			opts = append(opts, rabbitmq.DurableQueue())
		}

		if options.RequeueOnError {
			opts = append(opts, rabbitmq.RequeueOnError())
		}

		subOpts = broker.NewSubscribeOptions(opts...)
	case NATS:
		b = nats.NewBroker(bo...)
	default:
		b = memory.NewBroker(bo...)
	}

	return b, subOpts
}
