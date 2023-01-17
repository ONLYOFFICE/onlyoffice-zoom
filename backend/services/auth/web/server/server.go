package server

import (
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/rpc"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/adapter"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/service"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/handler"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/message"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	mclient "go-micro.dev/v4/client"
)

type AuthRPCServer struct {
	service port.UserAccessService
	zoomAPI client.ZoomAuth
	logger  log.Logger
}

func NewAuthRPCServer(opts ...Option) rpc.RPCEngine {
	options := newOptions(opts...)

	adptr := adapter.NewMemoryUserAdapter()
	if options.PersistenceURL != "" {
		adptr = adapter.NewMongoUserAdapter(options.PersistenceURL)
	}

	service := service.NewUserService(options.Logger, adptr, crypto.NewAesEncryptor([]byte(options.ClientSecret)))
	return AuthRPCServer{
		service: service,
		zoomAPI: client.NewZoomClient(options.ClientID, options.ClientSecret),
		logger:  options.Logger,
	}
}

func (a AuthRPCServer) BuildMessageHandlers() []rpc.RPCMessageHandler {
	return []rpc.RPCMessageHandler{
		{
			Topic:   "insert-auth",
			Queue:   "zoom-auth",
			Handler: message.BuildInsertMessageHandler(a.service).GetHandler(),
		},
		{
			Topic:   "delete-auth",
			Queue:   "zoom-auth",
			Handler: message.BuildDeleteMessageHandler(a.service).GetHandler(),
		},
	}
}

func (a AuthRPCServer) BuildHandlers(client mclient.Client) []interface{} {
	return []interface{}{
		handler.NewUserSelectHandler(a.service, client, a.zoomAPI, a.logger),
	}
}
