package cache

type CacheType int

var (
	Memory CacheType = 1
	Redis  CacheType = 2
	Micro  CacheType = 3
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	CacheType CacheType
	Size      int
	Address   string
	Password  string
	DB        int
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		Size: 100,
		DB:   0,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithCacheType sets cache type
func WithCacheType(val CacheType) Option {
	return func(o *Options) {
		o.CacheType = val
	}
}

// WithSize sets cache max size
func WithSize(val int) Option {
	return func(o *Options) {
		if val > 0 {
			o.Size = val
		}
	}
}

func WithAddress(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Address = val
		}
	}
}

func WithPassword(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Password = val
		}
	}
}

func WithDB(val int) Option {
	return func(o *Options) {
		if val >= 0 {
			o.DB = val
		}
	}
}
