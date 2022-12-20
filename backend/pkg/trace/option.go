package trace

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	Enable        bool
	Name          string
	Address       string
	TracerType    TracerType
	FractionRatio float64
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		Enable:        false,
		FractionRatio: 1,
		TracerType:    Default,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithEnable sets enable flag
func WithEnable(val bool) Option {
	return func(o *Options) {
		o.Enable = val
	}
}

// WithName sets service name to trace.
func WithName(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Name = val
		}
	}
}

// WithAddress sets remote distributed tracing address.
func WithAddress(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Address = val
		}
	}
}

// WithTracerType sets tracing type.
func WithTracerType(val TracerType) Option {
	return func(o *Options) {
		if val > 0 {
			o.TracerType = val
		}
	}
}

// WithFractionRation sets tracing fraction ratio.
func WithFractionRatio(val float64) Option {
	return func(o *Options) {
		if val > 0 {
			o.FractionRatio = val
		}
	}
}
