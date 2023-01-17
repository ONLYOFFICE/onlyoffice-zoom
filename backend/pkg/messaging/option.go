package messaging

type BrokerType int

var (
	RabbitMQ BrokerType = 1
	NATS     BrokerType = 2
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	BrokerType     BrokerType
	Addrs          []string
	DisableAutoAck bool
	Durable        bool
	AckOnSuccess   bool
	RequeueOnError bool
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		DisableAutoAck: false,
		Durable:        false,
		AckOnSuccess:   true,
		RequeueOnError: false,
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

func WithDisableAutoAck(val bool) Option {
	return func(o *Options) {
		o.DisableAutoAck = val
	}
}

func WithDurable(val bool) Option {
	return func(o *Options) {
		o.Durable = val
	}
}

func WithAckOnSuccess(val bool) Option {
	return func(o *Options) {
		o.AckOnSuccess = val
	}
}

func WithRequeueOnError(val bool) Option {
	return func(o *Options) {
		o.RequeueOnError = val
	}
}
