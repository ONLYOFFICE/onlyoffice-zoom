package worker

import (
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
)

type WorkerType int

var (
	Asynq WorkerType = 0
)

type WorkerOption func(*WorkerOptions)

type WorkerRedisCredentials struct {
	Addresses    []string
	Username     string
	Password     string
	Database     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Options defines the available options.
type WorkerOptions struct {
	MaxConcurrency   int
	RedisCredentials WorkerRedisCredentials
	WorkerType       WorkerType
	Logger           log.Logger
}

// NewOptions initializes the options.
func NewWorkerOptions(opts ...WorkerOption) WorkerOptions {
	opt := WorkerOptions{
		MaxConcurrency: 4,
		RedisCredentials: WorkerRedisCredentials{
			Addresses:    []string{"0.0.0.0:6379"},
			Database:     0,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		WorkerType: 0,
		Logger:     log.NewDefaultLogger(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithMaxConcurrency(val int) WorkerOption {
	return func(o *WorkerOptions) {
		if val > 0 {
			o.MaxConcurrency = val
		}
	}
}

func WithRedisCredentials(val WorkerRedisCredentials) WorkerOption {
	return func(o *WorkerOptions) {
		o.RedisCredentials = val
	}
}

func WithWorkerType(val WorkerType) WorkerOption {
	return func(wo *WorkerOptions) {
		wo.WorkerType = val
	}
}

func WithLogger(val log.Logger) WorkerOption {
	return func(o *WorkerOptions) {
		if val != nil {
			o.Logger = val
		}
	}
}

type EnqueuerOption func(*EnqueuerOptions)

type EnqueuerOptions struct {
	MaxRetry int
	Timeout  time.Duration
}

func NewEnqueuerOptions(opts ...EnqueuerOption) EnqueuerOptions {
	opt := EnqueuerOptions{
		MaxRetry: 3,
		Timeout:  0 * time.Second,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithMaxRetry(val int) EnqueuerOption {
	return func(eo *EnqueuerOptions) {
		if val > 0 {
			eo.MaxRetry = val
		}
	}
}

func WithTimeout(val time.Duration) EnqueuerOption {
	return func(eo *EnqueuerOptions) {
		eo.Timeout = val
	}
}
