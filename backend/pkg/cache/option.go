package cache

type CacheType int

var (
	Memory CacheType = 1
	Micro  CacheType = 2
)

// Option defines a single option.
type Option func(*Options)

// Options defines the available options.
type Options struct {
	CacheType CacheType
	Size      int
}

// NewOptions initializes the options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		Size: 100,
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
