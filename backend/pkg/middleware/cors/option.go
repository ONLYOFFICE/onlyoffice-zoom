package cors

// Option defines a single option.
type Option func(*Options)

// Options defines the set of available options.
type Options struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// NewOptions creates CORS middleware options.
func NewOptions(opts ...Option) Options {
	opt := Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithAllowedOrigins sets CORS AllowedOrigins header.
func WithAllowedOrigins(val ...string) Option {
	return func(o *Options) {
		if len(val) > 0 {
			o.AllowedOrigins = val
		}
	}
}

// WithAllowedMethods sets CORS AllowedMethods header.
func WithAllowedMethods(val ...string) Option {
	return func(o *Options) {
		if len(val) > 0 {
			o.AllowedMethods = val
		}
	}
}

// WithAllowedHeaders sets CORS AllowedHeaders header.
func WithAllowedHeaders(val ...string) Option {
	return func(o *Options) {
		if len(val) > 0 {
			o.AllowedHeaders = val
		}
	}
}

// WithAllowCredentials sets CORS AllowCredentials header.
func WithAllowCredentials(val bool) Option {
	return func(o *Options) {
		o.AllowCredentials = val
	}
}
