package server

import (
	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/rpc"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/adapter"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/service"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/handler"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/message"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	mclient "go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

type ConfigRPCServer struct {
	namespace   string
	zoomAPI     client.ZoomAuth
	service     port.SessionService
	jwtManager  crypto.JwtManager
	logger      plog.Logger
	callbackURL string
}

func NewConfigRPCServer(opts ...Option) rpc.RPCEngine {
	options := NewOptions(opts...)

	sessionAdapter, err := adapter.NewRedisSessionAdapter(
		adapter.WithBufferSize(options.Redis.BufferSize),
		adapter.WithRedisAddresses(options.Redis.RedisAddresses),
		adapter.WithRedisUsername(options.Redis.RedisUsername),
		adapter.WithRedisPassword(options.Redis.RedisPassword),
		adapter.WithLogger(options.Logger),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	jwtManager, err := crypto.NewOnlyofficeJwtManager(options.DocSecret)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return ConfigRPCServer{
		namespace:   options.Namespace,
		zoomAPI:     client.NewZoomClient(options.ClientID, options.ClientSecret),
		service:     service.NewSessionService(options.Logger, sessionAdapter),
		jwtManager:  jwtManager,
		logger:      options.Logger,
		callbackURL: options.CallbackURL,
	}
}

func (a ConfigRPCServer) BuildMessageHandlers() []rpc.RPCMessageHandler {
	return []rpc.RPCMessageHandler{
		{
			Topic:   "remove-session",
			Queue:   "zoom-builder",
			Handler: message.BuildRemoveSessionMessageHandler(a.logger, a.service).GetHandler(),
		},
		{
			Topic:   "remove-owner-session",
			Queue:   "zoom-builder",
			Handler: message.BuildOwnerRemoveSessionMessageHandler(a.logger, a.service).GetHandler(),
		},
	}
}

func (a ConfigRPCServer) BuildHandlers(c mclient.Client) []interface{} {
	return []interface{}{
		handler.NewConfigHandler(a.namespace, a.logger, c, a.zoomAPI, a.service, a.jwtManager, a.callbackURL),
		handler.NewSessionHandler(a.logger, a.service),
	}
}
