package server

import (
	"net/http"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/worker"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/server/controller"
	workerh "github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/server/worker"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/client"
)

type CallbackService struct {
	namespace     string
	maxSize       int64
	uploadTimeout int
	mux           *chi.Mux
	client        client.Client
	cache         cache.Cache
	jwtManager    crypto.JwtManager
	worker        worker.BackgroundWorker
	enqueuer      worker.BackgroundEnqueuer
	logger        log.Logger
}

// ApplyMiddleware useed to apply http server middlewares.
func (s CallbackService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewService initializes http server with options.
func NewServer(opts ...Option) CallbackService {
	gin.SetMode(gin.ReleaseMode)

	options := newOptions(opts...)
	wopts := []worker.WorkerOption{
		worker.WithMaxConcurrency(options.Worker.MaxConcurrency),
		worker.WithRedisCredentials(worker.WorkerRedisCredentials{
			Addresses: options.Worker.RedisCredentials.Addresses,
			Username:  options.Worker.RedisCredentials.Username,
			Password:  options.Worker.RedisCredentials.Password,
			Database:  options.Worker.RedisCredentials.Database,
		}),
	}

	jwtManager, _ := crypto.NewOnlyofficeJwtManager(options.DocSecret)

	service := CallbackService{
		namespace:     options.Namespace,
		maxSize:       options.MaxSize,
		uploadTimeout: options.UploadTimeout,
		mux:           chi.NewRouter(),
		jwtManager:    jwtManager,
		worker:        worker.NewBackgroundWorker(wopts...),
		enqueuer:      worker.NewBackgroundEnqueuer(wopts...),
		logger:        options.Logger,
	}

	return service
}

// NewHandler returns http server engine.
func (s CallbackService) NewHandler(client client.Client, cache cache.Cache) interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer(client, cache)
}

// InitializeServer sets all injected dependencies.
func (s *CallbackService) InitializeServer(c client.Client, cache cache.Cache) *chi.Mux {
	s.client = c
	s.cache = cache
	s.worker.Register("zoom-callback-upload", workerh.NewCallbackWorker(s.namespace, c, s.uploadTimeout, s.logger).UploadFile)
	s.InitializeRoutes()
	s.worker.Run()
	return s.mux
}

// InitializeRoutes builds all http routes.
func (s *CallbackService) InitializeRoutes() {
	callbackController := controller.NewCallbackController(s.namespace, s.maxSize, s.client, s.enqueuer, s.jwtManager, s.logger)
	s.mux.Group(func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)
		r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
			http.Redirect(rw, r.WithContext(r.Context()), "https://onlyoffice.com", http.StatusMovedPermanently)
		})
		r.Post("/callback", callbackController.BuildPostHandleCallback())
	})
}
