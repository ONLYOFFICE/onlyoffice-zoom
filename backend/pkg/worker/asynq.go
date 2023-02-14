package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

type asynqWorker struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

type asynqEnqueuer struct {
	client *asynq.Client
}

func newAsynqWorker(opts ...WorkerOption) BackgroundWorker {
	options := NewWorkerOptions(opts...)

	var workerOpts asynq.RedisConnOpt = asynq.RedisClientOpt{
		Addr:         options.RedisCredentials.Addresses[0],
		Username:     options.RedisCredentials.Username,
		Password:     options.RedisCredentials.Password,
		ReadTimeout:  options.RedisCredentials.ReadTimeout,
		WriteTimeout: options.RedisCredentials.WriteTimeout,
	}
	if len(options.RedisCredentials.Addresses) > 1 {
		workerOpts = asynq.RedisClusterClientOpt{
			Addrs:        options.RedisCredentials.Addresses,
			Username:     options.RedisCredentials.Username,
			Password:     options.RedisCredentials.Password,
			ReadTimeout:  options.RedisCredentials.ReadTimeout,
			WriteTimeout: options.RedisCredentials.WriteTimeout,
		}
	}

	return asynqWorker{
		srv: asynq.NewServer(workerOpts, asynq.Config{
			Concurrency: options.MaxConcurrency,
			Logger:      options.Logger,
		}),
		mux: asynq.NewServeMux(),
	}
}

func (w asynqWorker) Register(pattern string, handler func(ctx context.Context, payload []byte) error) {
	w.mux.Handle(pattern, asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		return handler(ctx, t.Payload())
	}))
}

func (w asynqWorker) Run() {
	go func() {
		if err := w.srv.Run(w.mux); err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func newAsynqEnqueuer(opts ...WorkerOption) BackgroundEnqueuer {
	options := NewWorkerOptions(opts...)

	var enqOpts asynq.RedisConnOpt = asynq.RedisClientOpt{
		Addr:         options.RedisCredentials.Addresses[0],
		Username:     options.RedisCredentials.Username,
		Password:     options.RedisCredentials.Password,
		ReadTimeout:  options.RedisCredentials.ReadTimeout,
		WriteTimeout: options.RedisCredentials.WriteTimeout,
	}
	if len(options.RedisCredentials.Addresses) > 1 {
		enqOpts = asynq.RedisClusterClientOpt{
			Addrs:        options.RedisCredentials.Addresses,
			Username:     options.RedisCredentials.Username,
			Password:     options.RedisCredentials.Password,
			ReadTimeout:  options.RedisCredentials.ReadTimeout,
			WriteTimeout: options.RedisCredentials.WriteTimeout,
		}
	}

	return asynqEnqueuer{
		client: asynq.NewClient(enqOpts),
	}
}

func (e asynqEnqueuer) Enqueue(pattern string, task []byte, opts ...EnqueuerOption) error {
	options := NewEnqueuerOptions(opts...)
	t := asynq.NewTask(pattern, task)

	_, err := e.client.Enqueue(t, asynq.MaxRetry(options.MaxRetry), asynq.Timeout(options.Timeout))
	return err
}

func (e asynqEnqueuer) EnqueueContext(ctx context.Context, pattern string, task []byte, opts ...EnqueuerOption) error {
	options := NewEnqueuerOptions(opts...)
	t := asynq.NewTask(pattern, task)

	_, err := e.client.EnqueueContext(ctx, t, asynq.MaxRetry(options.MaxRetry), asynq.Timeout(options.Timeout))
	return err
}

func (e asynqEnqueuer) Close() error {
	return e.Close()
}
