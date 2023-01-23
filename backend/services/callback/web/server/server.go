package server

import (
	"net/http"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/worker"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/server/controller"
	sworker "github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/server/worker"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gocraft/work"
	"go-micro.dev/v4/client"
)

type CallbackService struct {
	namespace     string
	mux           *chi.Mux
	client        client.Client
	logger        log.Logger
	jwtManager    crypto.JwtManager
	worker        *work.WorkerPool
	enqueuer      *work.Enqueuer
	maxSize       int64
	uploadTimeout int
}

// ApplyMiddleware useed to apply http server middlewares.
func (s CallbackService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewService initializes http server with options.
func NewServer(opts ...Option) CallbackService {
	gin.SetMode(gin.ReleaseMode)

	options := newOptions(opts...)
	wopts := []worker.Option{
		worker.WithMaxActive(options.Worker.MaxActive),
		worker.WithMaxIdle(options.Worker.MaxIdle),
		worker.WithRedisNamespace(options.Worker.RedisNamespace),
		worker.WithRedisAddress(options.Worker.RedisAddress),
		worker.WithRedisUsername(options.Worker.RedisUsername),
		worker.WithRedisPassword(options.Worker.RedisPassword),
		worker.WithRedisDatabase(options.Worker.RedisDatabase),
	}

	jwtManager, _ := crypto.NewOnlyofficeJwtManager(options.DocSecret)

	service := CallbackService{
		namespace:     options.Namespace,
		mux:           chi.NewRouter(),
		logger:        options.Logger,
		jwtManager:    jwtManager,
		worker:        worker.NewRedisWorker(sworker.NewWorkerContext(), wopts...),
		enqueuer:      worker.NewRedisEnqueuer(wopts...),
		maxSize:       options.MaxSize,
		uploadTimeout: options.UploadTimeout,
	}

	return service
}

// NewHandler returns http server engine.
func (s CallbackService) NewHandler(client client.Client) interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer(client)
}

// InitializeServer sets all injected dependencies.
func (s *CallbackService) InitializeServer(c client.Client) *chi.Mux {
	s.client = c
	s.worker.JobWithOptions("callback-upload", work.JobOptions{
		MaxFails: 3, SkipDead: true,
	}, sworker.NewCallbackWorker(s.namespace, c, s.uploadTimeout, s.logger).UploadFile)
	s.InitializeRoutes()
	s.worker.Start()
	return s.mux
}

// InitializeRoutes builds all http routes.
func (s *CallbackService) InitializeRoutes() {
	callbackController := controller.NewCallbackController(s.namespace, s.maxSize, s.logger, s.client, s.jwtManager)
	s.mux.Group(func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)
		r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
			http.Redirect(rw, r.WithContext(r.Context()), "https://onlyoffice.com", http.StatusMovedPermanently)
		})
		r.Post("/callback", callbackController.BuildPostHandleCallback(s.enqueuer))
	})
}
