package log

import (
	"os"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/shared"
)

type LogLevel int

const (
	LEVEL_TRACE   LogLevel = 1
	LEVEL_DEBUG   LogLevel = 2
	LEVEL_INFO    LogLevel = 3
	LEVEL_WARNING LogLevel = 4
	LEVEL_ERROR   LogLevel = 5
	LEVEL_FATAL   LogLevel = 6
)

// Option defines a single option.
type Option func(*Options)

// ElasticOption defines available Elastic options.
type ElasticOption struct {
	Address            string
	Index              string
	Bulk               bool
	Async              bool
	HealthcheckEnabled bool
	BasicAuthUsername  string
	BasicAuthPassword  string
	GzipEnabled        bool
}

// LumberjackOption defines available Lumberjack options
type LumberjackOption struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

// Options define the set of available options.
type Options struct {
	Name         string
	Machine      string
	Environment  string
	Level        LogLevel
	Pretty       bool
	Color        bool
	ReportCaller bool
	File         LumberjackOption
	Elastic      ElasticOption
}

// NewOptions creates logger options.
func NewOptions(opts ...Option) Options {
	host, _ := os.Hostname()

	if host == "" {
		host = "anonymous-machine"
	}

	options := Options{
		Name:         "anonymous-micro",
		Machine:      host,
		Environment:  shared.DEV,
		Level:        LEVEL_DEBUG,
		Pretty:       true,
		Color:        true,
		ReportCaller: false,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithName sets logger name.
func WithName(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Name = val
		}
	}
}

// WithMachine sets machine name.
func WithMachine(val string) Option {
	return func(o *Options) {
		if val != "" {
			o.Machine = val
		}
	}
}

// WithEnvironment sets logger environment.
func WithEnvironment(val string) Option {
	return func(o *Options) {
		if env, ok := shared.SUPPORTED_ENVIRONMENTS[val]; ok {
			o.Environment = env
			switch env {
			case shared.DEV:
				o.Level = LEVEL_DEBUG
			case shared.TEST:
				o.Level = LEVEL_INFO
			case shared.PROD:
				o.Level = LEVEL_WARNING
			}
		}
	}
}

// WithLevel sets logger level.
func WithLevel(val LogLevel) Option {
	return func(o *Options) {
		if val > 0 && val < 7 {
			o.Level = val
		}
	}
}

// WithPretty sets logrus pretty flag.
func WithPretty(val bool) Option {
	return func(o *Options) {
		o.Pretty = val
	}
}

// WithColor sets logrus color flag.
func WithColor(val bool) Option {
	return func(o *Options) {
		o.Color = val
	}
}

// WithReportCaller sets logrus caller flag.
func WithReportCaller(val bool) Option {
	return func(o *Options) {
		o.ReportCaller = val
	}
}

// WithFile sets lumberjack options.
func WithFile(val LumberjackOption) Option {
	return func(o *Options) {
		if val.Filename != "" {
			o.File = val
		}
	}
}

// WithElastic sets elastic options.
func WithElastic(val ElasticOption) Option {
	return func(o *Options) {
		if val.Address != "" && val.Index != "" {
			o.Elastic = val
		}
	}
}
