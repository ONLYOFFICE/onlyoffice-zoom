package worker

import (
	"context"
)

type BackgroundEnqueuer interface {
	Enqueue(pattern string, task []byte, opts ...EnqueuerOption) error
	EnqueueContext(ctx context.Context, pattern string, task []byte, opts ...EnqueuerOption) error
	Close() error
}

func NewBackgroundEnqueuer(opts ...WorkerOption) BackgroundEnqueuer {
	options := NewWorkerOptions(opts...)

	switch options.WorkerType {
	case Asynq:
		return newAsynqEnqueuer(opts...)
	default:
		return newAsynqEnqueuer(opts...)
	}
}
