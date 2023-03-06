package server

import (
	"net/http"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/controller"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/ws"
	zoomAPI "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/olahol/melody"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/client"
)

type ZoomHTTPService struct {
	namespace      string
	clientID       string
	clientSecret   string
	webhookSecret  string
	redirectURI    string
	hystrixTimeout int
	mux            *chi.Mux
	ws             *melody.Melody
	client         client.Client
	cache          cache.Cache
	store          sessions.Store
	logger         log.Logger
}

// NewService initializes http server with options.
func NewServer(opts ...Option) ZoomHTTPService {
	options := newOptions(opts...)

	gin.SetMode(gin.ReleaseMode)

	service := ZoomHTTPService{
		namespace:      options.Namespace,
		clientID:       options.ClientID,
		clientSecret:   options.ClientSecret,
		webhookSecret:  options.WebhookSecret,
		redirectURI:    options.RedirectURI,
		hystrixTimeout: options.HystrixTimout,
		mux:            chi.NewRouter(),
		store:          sessions.NewCookieStore([]byte(options.ClientSecret)),
		logger:         options.Logger,
	}

	return service
}

// ApplyMiddleware useed to apply http server middlewares.
func (s ZoomHTTPService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewHandler returns http server engine.
func (s ZoomHTTPService) NewHandler(client client.Client, cache cache.Cache) interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer(client, cache)
}

// InitializeServer sets all injected dependencies.
func (s *ZoomHTTPService) InitializeServer(c client.Client, cache cache.Cache) *chi.Mux {
	s.client = c
	s.cache = cache
	s.InitializeRoutes()
	s.ws = melody.New()
	s.ws.HandleConnect(ws.NewOnConnectHandler(s.namespace, s.clientSecret, s.ws, s.client))
	s.client.Options().Broker.Subscribe("notify-session", ws.NewNotifyOnMessage(s.logger, s.ws))
	return s.mux
}

// InitializeRoutes builds all http routes.
func (s *ZoomHTTPService) InitializeRoutes() {
	ctxMiddleware := middleware.BuildHandleZoomContextMiddleware(s.logger, s.clientSecret)
	eventMiddleware := middleware.BuildHandleZoomEventMiddleware(s.logger, s.webhookSecret)

	installController := controller.NewInstallController(s.clientID, s.store, s.logger)
	authController := controller.NewAuthController(
		s.namespace, s.hystrixTimeout, s.store, s.client,
		zoomAPI.NewZoomClient(s.clientID, s.clientSecret), s.logger,
	)
	apiController := controller.NewAPIController(
		s.namespace, s.hystrixTimeout, s.client, s.cache, zoomAPI.NewZoomApiClient(), s.logger,
	)

	s.mux.Group(func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)

		r.NotFound(func(rw http.ResponseWriter, cr *http.Request) {
			http.Redirect(rw, cr.WithContext(cr.Context()), "/oauth/install", http.StatusMovedPermanently)
		})

		r.Route("/oauth", func(cr chi.Router) {
			cr.Use(chimiddleware.NoCache)
			cr.Get("/install", installController.BuildGetInstall(s.redirectURI))
			cr.Get("/auth", authController.BuildGetAuth(s.redirectURI))
		})

		r.Route("/api", func(cr chi.Router) {
			cr.Use(ctxMiddleware)
			cr.Get("/files", apiController.BuildGetFiles())
			cr.Get("/config", apiController.BuildGetConfig())
			cr.Delete("/session", apiController.BuildDeleteSession())
		})

		r.Get("/track/{mid}", func(w http.ResponseWriter, r *http.Request) {
			s.ws.HandleRequest(w, r)
		})

		r.Route("/deauth", func(cr chi.Router) {
			cr.Use(chimiddleware.NoCache, eventMiddleware)
			cr.Post("/", authController.BuildPostDeauth())
		})
	})
}
