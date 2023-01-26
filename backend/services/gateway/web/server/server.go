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
	"go-micro.dev/v4/client"
)

type ZoomHTTPService struct {
	namespace     string
	mux           *chi.Mux
	ws            *melody.Melody
	client        client.Client
	logger        log.Logger
	store         sessions.Store
	clientID      string
	clientSecret  string
	webhookSecret string
	redirectURI   string
}

// NewService initializes http server with options.
func NewServer(opts ...Option) ZoomHTTPService {
	options := newOptions(opts...)

	gin.SetMode(gin.ReleaseMode)

	service := ZoomHTTPService{
		namespace:     options.Namespace,
		mux:           chi.NewRouter(),
		logger:        options.Logger,
		clientID:      options.ClientID,
		clientSecret:  options.ClientSecret,
		webhookSecret: options.WebhookSecret,
		redirectURI:   options.RedirectURI,
		store:         sessions.NewCookieStore([]byte(options.ClientSecret)),
	}

	return service
}

// ApplyMiddleware useed to apply http server middlewares.
func (s ZoomHTTPService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewHandler returns http server engine.
func (s ZoomHTTPService) NewHandler(client client.Client) interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer(client)
}

// InitializeServer sets all injected dependencies.
func (s *ZoomHTTPService) InitializeServer(c client.Client) *chi.Mux {
	s.client = c
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

	installController := controller.NewInstallController(s.logger, s.store, s.clientID)
	authController := controller.NewAuthController(s.logger, s.store, s.client, zoomAPI.NewZoomClient(s.clientID, s.clientSecret))
	apiController := controller.NewAPIController(s.namespace, s.logger, s.client, zoomAPI.NewZoomApiClient())

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
			cr.Get("/me", apiController.BuildGetMe())
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
