package message

import (
	"context"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
)

type RemoveSessionMessageHandler struct {
	logger  plog.Logger
	service port.SessionService
}

func BuildRemoveSessionMessageHandler(logger plog.Logger, service port.SessionService) RemoveSessionMessageHandler {
	return RemoveSessionMessageHandler{
		logger:  logger,
		service: service,
	}
}

func (i RemoveSessionMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		if mid, ok := payload.(string); !ok {
			return _ErrInvalidHandlerPayload
		} else {
			return i.service.DeleteSession(ctx, mid)
		}
	}
}
