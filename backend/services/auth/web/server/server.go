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

func NewAuthRPCServer(
	logger log.Logger,
	persistenceURL,
	clientID,
	clientSecret string,
) rpc.RPCEngine {
	adptr := adapter.NewMemoryUserAdapter()

	if persistenceURL != "" {
		adptr = adapter.NewMongoUserAdapter(persistenceURL)
	}

	service := service.NewUserService(logger, adptr, crypto.NewAesEncryptor([]byte(clientSecret)))
	return AuthRPCServer{
		service: service,
		zoomAPI: client.NewZoomClient(clientID, clientSecret),
		logger:  logger,
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
		handler.NewUserHandler(a.service, client, a.zoomAPI, a.logger),
	}
}
