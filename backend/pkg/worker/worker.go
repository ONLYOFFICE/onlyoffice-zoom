package worker

import "context"

type BackgroundWorker interface {
	Register(pattern string, handler func(ctx context.Context, payload []byte) error)
	Run()
}

func NewBackgroundWorker(opts ...WorkerOption) BackgroundWorker {
	options := NewWorkerOptions(opts...)

	switch options.WorkerType {
	case Asynq:
		return newAsynqWorker(opts...)
	default:
		return newAsynqWorker(opts...)
	}
}
